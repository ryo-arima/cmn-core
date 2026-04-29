package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clientshare "github.com/ryo-arima/cmn-core/pkg/client/share"
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
	manager := clientshare.NewManager(*conf, profile)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		controller.SetOutputFormat(outputFormat)
	}

	// usecases
	userUC := usecase.NewUser(*conf, manager)
	groupUC := usecase.NewGroup(*conf, manager)
	memberUC := usecase.NewMember(*conf, manager)
	resourceUC := usecase.NewResourceUC(*conf, manager)

	// login — explicit SSO login (also triggered automatically on first command)
	rootCmd.AddCommand(controller.InitSSOLoginCmd(manager))

	// get <target>
	getCmd := &cobra.Command{Use: "get", Short: "get a resource"}
	getCmd.AddCommand(controller.NewUserGetCmd(userUC))
	getCmd.AddCommand(controller.NewGroupGetCmd(groupUC))
	getCmd.AddCommand(controller.NewResourceGetCmd(resourceUC))
	rootCmd.AddCommand(getCmd)

	// list <target>
	listCmd := &cobra.Command{Use: "list", Short: "list resources"}
	listCmd.AddCommand(controller.NewUserListCmd(userUC))
	listCmd.AddCommand(controller.NewGroupListCmd(groupUC))
	listCmd.AddCommand(controller.NewMemberListCmd(memberUC))
	listCmd.AddCommand(controller.NewResourceListCmd(resourceUC))
	listCmd.AddCommand(controller.NewResourceGroupListCmd(resourceUC))
	rootCmd.AddCommand(listCmd)

	// create <target>
	createCmd := &cobra.Command{Use: "create", Short: "create a resource"}
	createCmd.AddCommand(controller.NewGroupCreateCmd(groupUC))
	createCmd.AddCommand(controller.NewResourceCreateCmd(resourceUC))
	rootCmd.AddCommand(createCmd)

	// update <target>
	updateCmd := &cobra.Command{Use: "update", Short: "update a resource"}
	updateCmd.AddCommand(controller.NewUserUpdateCmd(userUC))
	updateCmd.AddCommand(controller.NewGroupUpdateCmd(groupUC))
	updateCmd.AddCommand(controller.NewResourceUpdateCmd(resourceUC))
	rootCmd.AddCommand(updateCmd)

	// delete <target>
	deleteCmd := &cobra.Command{Use: "delete", Short: "delete a resource"}
	deleteCmd.AddCommand(controller.NewGroupDeleteCmd(groupUC))
	deleteCmd.AddCommand(controller.NewResourceDeleteCmd(resourceUC))
	rootCmd.AddCommand(deleteCmd)

	// add <target>
	addCmd := &cobra.Command{Use: "add", Short: "add a member"}
	addCmd.AddCommand(controller.NewMemberAddCmd(memberUC))
	rootCmd.AddCommand(addCmd)

	// remove <target>
	removeCmd := &cobra.Command{Use: "remove", Short: "remove a resource"}
	removeCmd.AddCommand(controller.NewMemberRemoveCmd(memberUC))
	removeCmd.AddCommand(controller.NewResourceGroupRemoveCmd(resourceUC))
	rootCmd.AddCommand(removeCmd)

	// set <target>
	setCmd := &cobra.Command{Use: "set", Short: "set a resource group role"}
	setCmd.AddCommand(controller.NewResourceGroupSetCmd(resourceUC))
	rootCmd.AddCommand(setCmd)

	// show <target>
	showCmd := &cobra.Command{Use: "show", Short: "show resource details"}
	showCmd.AddCommand(controller.NewResourceShowCmd(resourceUC))
	rootCmd.AddCommand(showCmd)

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
