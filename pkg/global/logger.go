package global

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// GlobalLoggerConfig represents global logger configuration
type GlobalLoggerConfig struct {
	Component    string `json:"component" yaml:"component"`
	Service      string `json:"service" yaml:"service"`
	Level        string `json:"level" yaml:"level"`
	Structured   bool   `json:"structured" yaml:"structured"`
	EnableCaller bool   `json:"enable_caller" yaml:"enable_caller"`
	Output       string `json:"output" yaml:"output"`
}

// GlobalLoggerInterface defines the global logging interface
type GlobalLoggerInterface interface {
	DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	INFO(mcode MCode, optionalMessage string)
	WARN(mcode MCode, optionalMessage string)
	ERROR(mcode MCode, optionalMessage string)
	FATAL(mcode MCode, optionalMessage string)
}

// GlobalLogLevel represents the log level
type GlobalLogLevel int

const (
	GLOBAL_DEBUG GlobalLogLevel = iota
	GLOBAL_INFO
	GLOBAL_WARN
	GLOBAL_ERROR
	GLOBAL_FATAL
)

// String returns string representation of log level
func (l GlobalLogLevel) String() string {
	switch l {
	case GLOBAL_DEBUG:
		return "DEBUG  "
	case GLOBAL_INFO:
		return "INFO   "
	case GLOBAL_WARN:
		return "WARN   "
	case GLOBAL_ERROR:
		return "ERROR  "
	case GLOBAL_FATAL:
		return "FATAL  "
	default:
		return "UNKNOWN"
	}
}

// GlobalLogEntry represents a structured log entry
type GlobalLogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Code      string                 `json:"code"`
	Component string                 `json:"component"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Line      int                    `json:"line,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// GlobalLogger represents the global application logger
type GlobalLogger struct {
	config *GlobalLoggerConfig
	level  GlobalLogLevel
	output io.Writer
}

// NewGlobalLogger creates a new global logger instance
func NewGlobalLogger(config GlobalLoggerConfig) GlobalLoggerInterface {
	logger := &GlobalLogger{
		config: &config,
		output: os.Stdout,
	}

	// Set log level
	switch strings.ToUpper(config.Level) {
	case "DEBUG":
		logger.level = GLOBAL_DEBUG
	case "INFO":
		logger.level = GLOBAL_INFO
	case "WARN":
		logger.level = GLOBAL_WARN
	case "ERROR":
		logger.level = GLOBAL_ERROR
	case "FATAL":
		logger.level = GLOBAL_FATAL
	default:
		logger.level = GLOBAL_INFO
	}

	// Set output
	switch config.Output {
	case "stderr":
		logger.output = os.Stderr
	case "stdout", "":
		logger.output = os.Stdout
	default:
		// File output
		if file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666); err == nil {
			logger.output = file
		} else {
			logger.output = os.Stdout
			fmt.Fprintf(os.Stderr, "Failed to open log file %s: %v\n", config.Output, err)
		}
	}

	return logger
}

// formatWithOptional formats the message with optional additional message
func formatGlobalWithOptional(mcode MCode, optionalMessage string) string {
	if optionalMessage == "" {
		return mcode.Message
	}
	return fmt.Sprintf("%s: %s", mcode.Message, optionalMessage)
}

// log writes a log entry
func (l *GlobalLogger) log(level GlobalLogLevel, mcode MCode, optionalMessage string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	finalMessage := formatGlobalWithOptional(mcode, optionalMessage)

	now := time.Now().UTC()
	timestamp := now.Format("2006-01-02T15:04:05.000000000") + " UTC"

	entry := GlobalLogEntry{
		Timestamp: timestamp,
		Level:     level.String(),
		Code:      mcode.Code,
		Component: l.config.Component,
		Service:   l.config.Service,
		Message:   finalMessage,
		Fields:    fields,
	}

	l.writeLogEntry(entry)
}

// writeLogEntry writes the actual log entry to output
func (l *GlobalLogger) writeLogEntry(entry GlobalLogEntry) {
	// Add caller information if enabled or DEBUG level
	if l.config.EnableCaller || l.level == GLOBAL_DEBUG {
		if pc, file, line, ok := runtime.Caller(4); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	// Extract error field
	if entry.Fields != nil {
		if err, ok := entry.Fields["error"].(string); ok {
			entry.Error = err
			delete(entry.Fields, "error")
		}
		if err, ok := entry.Fields["error"].(error); ok {
			entry.Error = err.Error()
			delete(entry.Fields, "error")
		}
	}

	if l.config.Structured {
		// JSON format
		if jsonBytes, err := json.Marshal(entry); err == nil {
			fmt.Fprintln(l.output, string(jsonBytes))
		} else {
			// Fallback
			fmt.Fprintf(l.output, "[%s] [%s] [%s] %s\n",
				entry.Timestamp, entry.Level, entry.Code, entry.Message)
		}
	} else {
		// Human-readable format
		fmt.Fprintf(l.output, "[%s] [%s] [%s] %s",
			entry.Timestamp, entry.Level, entry.Code, entry.Message)
		if len(entry.Fields) > 0 && entry.Level == "DEBUG  " {
			if fieldsJSON, err := json.Marshal(entry.Fields); err == nil {
				fmt.Fprintf(l.output, " %s", string(fieldsJSON))
			}
		}
		fmt.Fprintln(l.output)
	}
}

// DEBUG logs a debug message
func (l *GlobalLogger) DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(GLOBAL_DEBUG, mcode, optionalMessage, f)
}

// INFO logs an info message
func (l *GlobalLogger) INFO(mcode MCode, optionalMessage string) {
	l.log(GLOBAL_INFO, mcode, optionalMessage, nil)
}

// WARN logs a warning message
func (l *GlobalLogger) WARN(mcode MCode, optionalMessage string) {
	l.log(GLOBAL_WARN, mcode, optionalMessage, nil)
}

// ERROR logs an error message
func (l *GlobalLogger) ERROR(mcode MCode, optionalMessage string) {
	l.log(GLOBAL_ERROR, mcode, optionalMessage, nil)
}

// FATAL logs a fatal message and exits
func (l *GlobalLogger) FATAL(mcode MCode, optionalMessage string) {
	l.log(GLOBAL_FATAL, mcode, optionalMessage, nil)
	os.Exit(1)
}
