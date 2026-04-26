package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/controller"
	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cmn-app",
		Short: "cmn-core app client (accesses /v1/internal/* endpoints)",
		Long:  "Command-line interface for app-role operations. Authentication is handled automatically via the configured SSO provider.",
	}

	var configFile string
	var outputFormat string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "etc/app.yaml", "path to config file (env: CONFIG_FILE, default: etc/app.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table, json, yaml")

	// フラグを早期解析して config.NewBaseConfig() の前に CONFIG_FILE を設定する
	_ = rootCmd.ParseFlags(os.Args[1:])
	os.Setenv("CONFIG_FILE", configFile)

	conf := config.NewBaseConfig()
	profile := strings.TrimSuffix(filepath.Base(configFile), ".yaml")
	manager := auth.NewManager(*conf, profile)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}

	// idp / resource usecases
	idpUC := usecase.NewIdP(*conf, manager)
	resourceUC := usecase.NewResourceUC(*conf, manager)

	// login — explicit SSO login (also triggered automatically on first command)
	rootCmd.AddCommand(controller.InitSSOLoginCmd(manager))

	// user / group / member / resource subcommands
	rootCmd.AddCommand(controller.InitUserCmd(idpUC))
	rootCmd.AddCommand(controller.InitGroupCmd(idpUC))
	rootCmd.AddCommand(controller.InitMemberCmd(idpUC))
	rootCmd.AddCommand(controller.InitResourceCmd(resourceUC))

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
