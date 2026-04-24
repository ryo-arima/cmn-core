package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/client/controller"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	conf := config.NewBaseConfig()

	rootCmd := &cobra.Command{
		Use:   "cmn",
		Short: "cmn-core anonymous client (accesses /v1/public/* endpoints)",
		Long:  "Command-line interface for unauthenticated public operations. Use --access-token or CMN_ACCESS_TOKEN to validate a token.",
	}

	var outputFormat string
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table, json, yaml")

	// token subcommands
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "token operations",
	}
	tokenCmd.AddCommand(controller.InitAnonymousValidateCmd(*conf))
	rootCmd.AddCommand(tokenCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := os.ReadFile("VERSION")
			fmt.Printf("cmn %s\n", strings.TrimSpace(string(b)))
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
