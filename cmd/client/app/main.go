package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/controller"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	conf := config.NewBaseConfig()
	manager := auth.NewManager(*conf, "app")

	rootCmd := &cobra.Command{
		Use:   "cmn-app",
		Short: "cmn-core app client (accesses /v1/internal/* endpoints)",
		Long:  "Command-line interface for app-role operations. Authentication is handled automatically via the configured SSO provider.",
	}

	var outputFormat string
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table, json, yaml")

	// login — explicit SSO login (also triggered automatically on first command)
	rootCmd.AddCommand(controller.InitSSOLoginCmd(manager))

	// token subcommands
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "manage authentication tokens",
	}
	tokenCmd.AddCommand(controller.InitCommonRefreshTokenCmd(manager))
	tokenCmd.AddCommand(controller.InitCommonLogoutCmd(manager))
	tokenCmd.AddCommand(controller.InitCommonValidateTokenCmd(manager))
	tokenCmd.AddCommand(controller.InitCommonUserInfoCmd(manager))
	rootCmd.AddCommand(tokenCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := os.ReadFile("VERSION")
			fmt.Printf("cmn-app %s\n", strings.TrimSpace(string(b)))
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
