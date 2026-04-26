package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/controller"
	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cmn-admin",
		Short: "cmn-core admin client (accesses /v1/private/* endpoints)",
		Long:  "Command-line interface for admin-role operations. Authentication is handled automatically via the configured SSO provider.",
	}

	var configFile string
	var outputFormat string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to config file (env: CONFIG_FILE, default: etc/app.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format: table, json, yaml")

	// フラグを早期解析して config.NewBaseConfig() の前に CONFIG_FILE を設定する
	_ = rootCmd.ParseFlags(os.Args[1:])
	if configFile != "" {
		os.Setenv("CONFIG_FILE", configFile)
	}

	conf := config.NewBaseConfig()
	manager := auth.NewManager(*conf, "admin")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}

	// idp-admin / resource-admin usecases
	idpAdminUC := usecase.NewIdPAdmin(*conf, manager)
	resourceAdminUC := usecase.NewResourceAdminUC(*conf, manager)

	// login — explicit SSO login (also triggered automatically on first command)
	rootCmd.AddCommand(controller.InitSSOLoginCmd(manager))

	// user / group / member / resource subcommands (admin)
	rootCmd.AddCommand(controller.InitAdminUserCmd(idpAdminUC))
	rootCmd.AddCommand(controller.InitAdminGroupCmd(idpAdminUC))
	rootCmd.AddCommand(controller.InitAdminMemberCmd(idpAdminUC))
	rootCmd.AddCommand(controller.InitResourceCmd(resourceAdminUC))

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
			fmt.Printf("cmn-admin %s\n", strings.TrimSpace(string(b)))
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
