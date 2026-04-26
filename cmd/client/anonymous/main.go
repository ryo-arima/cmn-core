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
	rootCmd := &cobra.Command{
		Use:   "cmn",
		Short: "cmn-core anonymous client (accesses /v1/public/* endpoints)",
		Long:  "Command-line interface for unauthenticated public operations. Use --access-token or CMN_ACCESS_TOKEN to validate a token.",
	}

	var configFile string
	var outputFormat string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "etc/app.yaml", "path to config file (env: CONFIG_FILE, default: etc/app.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table, json, yaml")

	// フラグを早期解析して config.NewBaseConfig() の前に CONFIG_FILE を設定する
	_ = rootCmd.ParseFlags(os.Args[1:])
	os.Setenv("CONFIG_FILE", configFile)

	conf := config.NewBaseConfig()

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}

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
