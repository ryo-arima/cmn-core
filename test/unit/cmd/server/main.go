package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rodaine/table"
)

// stripANSI removes ANSI color codes from a string for width calculation
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}

// ansiAwareWidth calculates string width ignoring ANSI color codes
func ansiAwareWidth(s string) int {
	return len([]rune(stripANSI(s)))
}

// Styles
var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("57")).
		Padding(0, 1).
		MarginBottom(1)
)

type TestResult struct {
	Type    string
	Message string
	Package string
	Test    string
}

type TestTableRow struct {
	No       string
	MainTest string
	SubTest  string
	Status   string
	Results  string
	Duration string
}

func parseTestLine(line string) TestResult {
	line = strings.TrimSpace(line)

	// Regular expressions for different test output patterns
	runPattern := regexp.MustCompile(`^=== RUN\s+(.+)`)
	passPattern := regexp.MustCompile(`^--- PASS:\s+(.+)\s+\((.+)\)`)
	failPattern := regexp.MustCompile(`^--- FAIL:\s+(.+)\s+\((.+)\)`)
	pkgPassPattern := regexp.MustCompile(`^ok\s+(.+)\s+(.+)`)
	pkgFailPattern := regexp.MustCompile(`^FAIL\s+(.+)\s+(.+)`)
	coveragePattern := regexp.MustCompile(`coverage:\s+(.+)`)

	switch {
	case runPattern.MatchString(line):
		matches := runPattern.FindStringSubmatch(line)
		return TestResult{Type: "RUN", Message: matches[1]}
	case passPattern.MatchString(line):
		matches := passPattern.FindStringSubmatch(line)
		return TestResult{Type: "PASS", Message: fmt.Sprintf("%s (%s)", matches[1], matches[2])}
	case failPattern.MatchString(line):
		matches := failPattern.FindStringSubmatch(line)
		return TestResult{Type: "FAIL", Message: fmt.Sprintf("%s (%s)", matches[1], matches[2])}
	case pkgPassPattern.MatchString(line):
		matches := pkgPassPattern.FindStringSubmatch(line)
		return TestResult{Type: "PASS", Package: matches[1], Message: fmt.Sprintf("Package completed (%s)", matches[2])}
	case pkgFailPattern.MatchString(line):
		matches := pkgFailPattern.FindStringSubmatch(line)
		return TestResult{Type: "FAIL", Package: matches[1], Message: fmt.Sprintf("Package failed (%s)", matches[2])}
	case coveragePattern.MatchString(line):
		matches := coveragePattern.FindStringSubmatch(line)
		return TestResult{Type: "COVERAGE", Message: matches[1]}
	default:
		if strings.Contains(line, "panic:") {
			return TestResult{Type: "ERROR", Message: line}
		}
		if line != "" {
			return TestResult{Type: "INFO", Message: line}
		}
	}

	return TestResult{}
}

// CoverageInfo holds coverage information for a function
type CoverageInfo struct {
	Function   string
	Percent    float64
	PercentStr string
}

// CoverageStats holds statistical coverage information
type CoverageStats struct {
	Category   string
	Percent    float64
	PercentStr string
}

// parseDetailedCoverage parses detailed coverage information from coverage output
func parseDetailedCoverage() []CoverageStats {
	var coverageStats []CoverageStats

	// Read coverage.out file directly to analyze lines and statements
	file, err := os.Open("test/results/coverage.out")
	if err != nil {
		return coverageStats
	}
	defer file.Close()

	var totalStatements int
	var coveredStatements int
	var totalLines int
	var coveredLines int

	// Parse coverage.out file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		// Parse lines like:
		// github.com/ryo-arima/cmn-core/test/pkg/server/repository/test_helper.go:25.34,61.2 3 1
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			// Extract statement count and coverage info
			if statementCount, err := strconv.Atoi(parts[1]); err == nil {
				totalStatements += statementCount

				// Check if covered (last field is 1 for covered, 0 for not covered)
				if covered, err := strconv.Atoi(parts[2]); err == nil && covered > 0 {
					coveredStatements += statementCount
				}
			}

			// Extract line range information
			locationPart := parts[0]
			if colonIndex := strings.LastIndex(locationPart, ":"); colonIndex != -1 {
				rangePart := locationPart[colonIndex+1:]
				if commaIndex := strings.Index(rangePart, ","); commaIndex != -1 {
					startPart := rangePart[:commaIndex]
					endPart := rangePart[commaIndex+1:]

					// Extract start and end line numbers
					if dotIndex := strings.Index(startPart, "."); dotIndex != -1 {
						if startLine, err := strconv.Atoi(startPart[:dotIndex]); err == nil {
							if dotIndex := strings.Index(endPart, "."); dotIndex != -1 {
								if endLine, err := strconv.Atoi(endPart[:dotIndex]); err == nil {
									lineCount := endLine - startLine + 1
									totalLines += lineCount

									// Check if this block is covered
									if covered, err := strconv.Atoi(parts[2]); err == nil && covered > 0 {
										coveredLines += lineCount
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Calculate coverage percentages
	if totalStatements > 0 {
		statementCoverage := float64(coveredStatements) / float64(totalStatements) * 100
		coverageStats = append(coverageStats, CoverageStats{
			Category:   "Statements Covered",
			Percent:    statementCoverage,
			PercentStr: fmt.Sprintf("%.1f%% (%d/%d)", statementCoverage, coveredStatements, totalStatements),
		})
	}

	if totalLines > 0 {
		lineCoverage := float64(coveredLines) / float64(totalLines) * 100
		coverageStats = append(coverageStats, CoverageStats{
			Category:   "Lines Covered",
			Percent:    lineCoverage,
			PercentStr: fmt.Sprintf("%.1f%% (%d/%d)", lineCoverage, coveredLines, totalLines),
		})
	}

	// Calculate branch coverage approximation from function coverage
	cmd := exec.Command("go", "tool", "cover", "-func=test/results/coverage.out")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		var totalFunctions int
		var coveredFunctions int

		for _, line := range lines {
			if line == "" || strings.Contains(line, "total:") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) >= 2 && strings.HasSuffix(parts[len(parts)-1], "%") {
				totalFunctions++
				percentStr := strings.TrimSuffix(parts[len(parts)-1], "%")
				if percent, err := strconv.ParseFloat(percentStr, 64); err == nil && percent > 0 {
					coveredFunctions++
				}
			}
		}

		if totalFunctions > 0 {
			branchCoverage := float64(coveredFunctions) / float64(totalFunctions) * 100
			coverageStats = append(coverageStats, CoverageStats{
				Category:   "Functions Covered",
				Percent:    branchCoverage,
				PercentStr: fmt.Sprintf("%.1f%% (%d/%d)", branchCoverage, coveredFunctions, totalFunctions),
			})
		}
	}

	return coverageStats
}

// extractCoveragePercent extracts coverage percentage from coverage output
func extractCoveragePercent() string {
	cmd := exec.Command("go", "tool", "cover", "-func=test/results/coverage.out")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "total:") {
			// Extract percentage from line like "total:		(statements)	67.8%"
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasSuffix(part, "%") {
					return part
				}
			}
		}
	}
	return "N/A"
}

// printTable prints a formatted table with proper alignment using rodaine/table
func printTable(rows []TestTableRow) {
	if len(rows) == 0 {
		return
	}

	// Create table with rodaine/table - separate columns for main and sub tests
	tbl := table.New("No.", "Main Test", "Sub Test", "Status", "Results", "Duration")

	// Set consistent padding
	tbl.WithPadding(1)

	// Use ANSI-aware width function to handle colored text properly
	tbl.WithWidthFunc(ansiAwareWidth)

	// Add rows to table
	for _, row := range rows {
		// Clean up the values to remove any existing styling for proper alignment
		no := strings.TrimSpace(row.No)
		mainTest := row.MainTest
		subTest := row.SubTest

		// Apply colors to Status and Results
		var status string
		if row.Status == "RUN" {
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Render(row.Status)
		} else {
			status = row.Status
		}

		var results string
		switch row.Results {
		case "PASS":
			results = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render(row.Results)
		case "FAIL":
			results = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(row.Results)
		default:
			results = row.Results
		}

		duration := row.Duration

		// Truncate test names if too long
		if len(mainTest) > 30 {
			mainTest = mainTest[:27] + "..."
		}
		if len(subTest) > 25 {
			subTest = subTest[:22] + "..."
		}

		tbl.AddRow(no, mainTest, subTest, status, results, duration)
	}

	// Print the table
	tbl.Print()
} // displayNonInteractiveStats displays beautiful test statistics in non-interactive mode
func displayNonInteractiveStats(tableRows []TestTableRow, testType string, failed bool) {
	fmt.Println()

	// Count statistics from table rows (which represent actual test execution)
	passed := 0
	totalFailed := 0
	totalTests := 0

	for _, row := range tableRows {
		// Only count rows that have actual test results (not just RUN status)
		if row.Results == "PASS" {
			passed++
			totalTests++
		} else if row.Results == "FAIL" {
			totalFailed++
			totalTests++
		}
		// Don't count RUN status or empty results in totals
	}

	// Determine colors based on test results
	var borderColor lipgloss.Color
	var statusColor lipgloss.Color

	if failed || totalFailed > 0 {
		borderColor = lipgloss.Color("196") // Red
		statusColor = lipgloss.Color("196") // Red
	} else {
		borderColor = lipgloss.Color("82") // Green
		statusColor = lipgloss.Color("82") // Green
	}

	// Create statistics display with dynamic colors
	statsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	var statsContent strings.Builder

	// Title
	var title string
	switch testType {
	case "repository":
		title = "Repository Test Results"
	case "controller":
		title = "Controller Test Results"
	case "coverage":
		title = "Coverage Test Results"
	case "server":
		title = "Server Test Results"
	default:
		title = fmt.Sprintf("%s Test Results", strings.ToUpper(testType))
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	statsContent.WriteString(titleStyle.Render(title))
	statsContent.WriteString("\n\n")

	// Test counts
	if totalTests > 0 {
		passedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
		failedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

		statsContent.WriteString(fmt.Sprintf("Total Tests: %d\n", totalTests))
		statsContent.WriteString(passedStyle.Render(fmt.Sprintf("Passed: %d", passed)))
		statsContent.WriteString("  ")

		if totalFailed > 0 {
			statsContent.WriteString(failedStyle.Render(fmt.Sprintf("Failed: %d", totalFailed)))
		} else {
			statsContent.WriteString("Failed: 0")
		}
		statsContent.WriteString("\n")
	}

	// Package count (fixed to 1 for single package tests)
	statsContent.WriteString("Packages: 1\n")

	// Status
	statusContent := "\n"
	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true)

	var statusText string
	if failed || totalFailed > 0 {
		statusText = "TESTS FAILED"
	} else {
		statusText = "ALL TESTS PASSED"
	}

	statusContent += statusStyle.Render(statusText)

	statsContent.WriteString(statusContent)

	// Display the styled statistics box
	fmt.Println(statsStyle.Render(statsContent.String()))
	fmt.Println()
}

// runNonInteractive runs tests in non-interactive mode with simple output
func runNonInteractive(testType string) {
	var testPath string

	switch testType {
	case "server":
		testPath = "./test/pkg/server/..."
	case "repository":
		testPath = "./test/pkg/server/repository/..."
	case "controller":
		testPath = "./test/pkg/server/controller/..."
	default:
		testPath = "./test/pkg/server/..."
	}

	// Header with styling
	var headerText string
	switch testType {
	case "repository":
		headerText = "Running Repository Tests with Coverage"
	case "controller":
		headerText = "Running Controller Tests with Coverage"
	case "server":
		headerText = "Running Server Tests with Coverage"
	default:
		headerText = fmt.Sprintf("Running %s Tests with Coverage", strings.ToUpper(testType))
	}

	fmt.Println(headerStyle.Render(headerText))
	fmt.Println()

	// Always run with coverage
	os.MkdirAll("test/results", 0755)
	cmd := exec.Command("go", "test", "-v", "-coverprofile=test/results/coverage.out", testPath)

	// Run command and capture output
	output, err := cmd.CombinedOutput()

	if err != nil && len(output) == 0 {
		fmt.Printf("Tests failed with no output: %v\n", err)
		os.Exit(1)
	}

	// Process and display output with formatting
	lines := strings.Split(string(output), "\n")
	var tableRows []TestTableRow
	testCounter := 0
	majorNumber := 1
	minorNumber := 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse test result for coloring and numbering
		result := parseTestLine(line)
		if result.Type != "" {

			// Add test numbering for RUN events
			if result.Type == "RUN" {
				testCounter++
				minorNumber++
				// Extract test name
				testName := strings.TrimPrefix(line, "=== RUN   ")

				// Check if this is a sub-test (contains '/')
				if strings.Contains(testName, "/") {
					// This is a sub-test, don't increment counters or add as main test
					testCounter--
					minorNumber--
					// Extract main test name and sub-test name
					parts := strings.SplitN(testName, "/", 2)
					if len(parts) == 2 {
						subTestName := parts[1]

						tableRows = append(tableRows, TestTableRow{
							No:       "",
							MainTest: "", // Empty for sub-tests to avoid duplication
							SubTest:  subTestName,
							Status:   "RUN",
							Results:  "-",
							Duration: "-",
						})
					}
				} else {
					// This is a main test
					// Format number as "1.1" with proper padding
					numberStr := fmt.Sprintf("%d.%d", majorNumber, minorNumber)

					tableRows = append(tableRows, TestTableRow{
						No:       numberStr,
						MainTest: testName,
						SubTest:  "",
						Status:   "RUN",
						Results:  "-",
						Duration: "-",
					})
				}
			} else if result.Type == "PASS" {
				// Update the matching row with PASS result
				if len(tableRows) > 0 {
					passLine := strings.TrimPrefix(line, "--- PASS: ")
					parts := strings.Split(passLine, " (")
					testName := parts[0]

					// Find the matching test row (search from end to beginning for recent tests)
					for i := len(tableRows) - 1; i >= 0; i-- {
						row := &tableRows[i]
						// Check if this PASS result matches either main test or sub test
						if (row.MainTest != "" && (strings.Contains(row.MainTest, testName) || strings.Contains(testName, row.MainTest))) ||
							(row.SubTest != "" && (strings.Contains(row.SubTest, testName) || strings.Contains(testName, row.SubTest))) ||
							(row.MainTest != "" && row.SubTest == "" && testName == row.MainTest) ||
							(row.SubTest != "" && strings.HasSuffix(testName, row.SubTest)) {

							duration := "-"
							if len(parts) > 1 {
								duration = "(" + parts[1]
								if len(duration) > 10 {
									duration = duration[:9] + ")"
								}
							}
							row.Results = "PASS"
							row.Duration = duration
							break
						}
					}
				}
			} else if result.Type == "FAIL" {
				// Update the matching row with FAIL result
				if len(tableRows) > 0 {
					failLine := strings.TrimPrefix(line, "--- FAIL: ")
					parts := strings.Split(failLine, " (")
					testName := parts[0]

					// Find the matching test row (search from end to beginning for recent tests)
					for i := len(tableRows) - 1; i >= 0; i-- {
						row := &tableRows[i]
						// Check if this FAIL result matches either main test or sub test
						if (row.MainTest != "" && (strings.Contains(row.MainTest, testName) || strings.Contains(testName, row.MainTest))) ||
							(row.SubTest != "" && (strings.Contains(row.SubTest, testName) || strings.Contains(testName, row.SubTest))) ||
							(row.MainTest != "" && row.SubTest == "" && testName == row.MainTest) ||
							(row.SubTest != "" && strings.HasSuffix(testName, row.SubTest)) {

							duration := "-"
							if len(parts) > 1 {
								duration = "(" + parts[1]
								if len(duration) > 10 {
									duration = duration[:9] + ")"
								}
							}
							row.Results = "FAIL"
							row.Duration = duration
							break
						}
					}
				}
			} else if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "FAIL ") {
				// Package completion - print current table and reset for next package
				if len(tableRows) > 0 {
					fmt.Println()
					printTable(tableRows)
				}

				fmt.Printf("\n%s\n", line)
				majorNumber++
				minorNumber = 0
				tableRows = []TestTableRow{} // Reset for next package
			} else if strings.Contains(line, "coverage:") {
				// Skip coverage lines - will be handled separately later
				continue
			}
			// Note: Sub-tests are now handled in the RUN section above
		}
	}

	// Print final table if there are remaining rows
	if len(tableRows) > 0 {
		fmt.Println()
		printTable(tableRows)
	}

	// Display coverage information using rodaine/table
	coveragePercent := extractCoveragePercent()
	if coveragePercent != "N/A" && coveragePercent != "" {
		fmt.Println()

		// Create coverage results table with detailed information
		coverageTbl := table.New("Coverage Category", "Percentage")
		coverageTbl.WithPadding(1)
		coverageTbl.WithWidthFunc(ansiAwareWidth)

		// Parse detailed coverage information first
		detailedCoverage := parseDetailedCoverage()
		for _, coverage := range detailedCoverage {
			var coloredPercent string
			if coverage.Percent >= 80 {
				coloredPercent = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render(coverage.PercentStr)
			} else if coverage.Percent >= 60 {
				coloredPercent = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(coverage.PercentStr)
			} else {
				coloredPercent = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(coverage.PercentStr)
			}
			coverageTbl.AddRow(coverage.Category, coloredPercent)
		}

		// Add total coverage last
		coverageTbl.AddRow("Total Coverage", lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(coveragePercent))

		coverageTbl.Print()
	}

	// Display beautiful statistics with coverage using actual table data
	displayNonInteractiveStats(tableRows, testType, false)
}

func main() {
	testType := "server"
	if len(os.Args) > 1 {
		testType = os.Args[1]
	}

	switch testType {
	case "server", "repository", "controller":
		// Valid test types
	default:
		fmt.Println("Usage: go run test/cmd/server/main.go [server|repository|controller]")
		fmt.Println("Note: Coverage is automatically included in all test runs")
		os.Exit(1)
	}

	// Always run in non-interactive mode
	runNonInteractive(testType)
}
