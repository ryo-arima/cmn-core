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
		Use:   "cmn-admin",
		Short: "cmn-core admin client (accesses /v1/private/* endpoints)",
		Long:  "Command-line interface for admin-role operations. Authentication is handled automatically via the configured SSO provider.",
	}

	var configFile string
	var outputFormat string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "etc/admin.yaml", "path to config file (env: CONFIG_FILE, default: etc/admin.yaml)")
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
	userAdminUC := usecase.NewUserAdmin(*conf, manager)
	groupAdminUC := usecase.NewGroupAdmin(*conf, manager)
	memberAdminUC := usecase.NewMemberAdmin(*conf, manager)
	resourceAdminUC := usecase.NewResourceAdminUC(*conf, manager)

	// login — explicit SSO login (also triggered automatically on first command)
	rootCmd.AddCommand(controller.InitSSOLoginCmd(manager))

	// get <target>
	getCmd := &cobra.Command{Use: "get", Short: "get a resource"}
	getCmd.AddCommand(controller.NewAdminUserGetCmd(userAdminUC))
	getCmd.AddCommand(controller.NewAdminGroupGetCmd(groupAdminUC))
	getCmd.AddCommand(controller.NewResourceGetCmd(resourceAdminUC))
	rootCmd.AddCommand(getCmd)

	// list <target>
	listCmd := &cobra.Command{Use: "list", Short: "list resources"}
	listCmd.AddCommand(controller.NewAdminUserListCmd(userAdminUC))
	listCmd.AddCommand(controller.NewAdminGroupListCmd(groupAdminUC))
	listCmd.AddCommand(controller.NewAdminMemberListCmd(memberAdminUC))
	listCmd.AddCommand(controller.NewResourceListCmd(resourceAdminUC))
	listCmd.AddCommand(controller.NewResourceGroupListCmd(resourceAdminUC))
	rootCmd.AddCommand(listCmd)

	// create <target>
	createCmd := &cobra.Command{Use: "create", Short: "create a resource"}
	createCmd.AddCommand(controller.NewAdminUserCreateCmd(userAdminUC))
	createCmd.AddCommand(controller.NewAdminGroupCreateCmd(groupAdminUC))
	createCmd.AddCommand(controller.NewResourceCreateCmd(resourceAdminUC))
	rootCmd.AddCommand(createCmd)

	// update <target>
	updateCmd := &cobra.Command{Use: "update", Short: "update a resource"}
	updateCmd.AddCommand(controller.NewAdminUserUpdateCmd(userAdminUC))
	updateCmd.AddCommand(controller.NewAdminGroupUpdateCmd(groupAdminUC))
	updateCmd.AddCommand(controller.NewResourceUpdateCmd(resourceAdminUC))
	rootCmd.AddCommand(updateCmd)

	// delete <target>
	deleteCmd := &cobra.Command{Use: "delete", Short: "delete a resource"}
	deleteCmd.AddCommand(controller.NewAdminUserDeleteCmd(userAdminUC))
	deleteCmd.AddCommand(controller.NewAdminGroupDeleteCmd(groupAdminUC))
	deleteCmd.AddCommand(controller.NewResourceDeleteCmd(resourceAdminUC))
	rootCmd.AddCommand(deleteCmd)

	// add <target>
	addCmd := &cobra.Command{Use: "add", Short: "add a member"}
	addCmd.AddCommand(controller.NewAdminMemberAddCmd(memberAdminUC))
	rootCmd.AddCommand(addCmd)

	// remove <target>
	removeCmd := &cobra.Command{Use: "remove", Short: "remove a resource"}
	removeCmd.AddCommand(controller.NewAdminMemberRemoveCmd(memberAdminUC))
	removeCmd.AddCommand(controller.NewResourceGroupRemoveCmd(resourceAdminUC))
	rootCmd.AddCommand(removeCmd)

	// set <target>
	setCmd := &cobra.Command{Use: "set", Short: "set a resource group role"}
	setCmd.AddCommand(controller.NewResourceGroupSetCmd(resourceAdminUC))
	rootCmd.AddCommand(setCmd)

	// show <target>
	showCmd := &cobra.Command{Use: "show", Short: "show resource details"}
	showCmd.AddCommand(controller.NewResourceShowCmd(resourceAdminUC))
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

	// bootstrap subcommands
	bootstrapCmd := &cobra.Command{Use: "bootstrap", Short: "bootstrap management"}
	bootstrapCmd.AddCommand(controller.NewDBBootstrapCmd(*conf))
	rootCmd.AddCommand(bootstrapCmd)

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
