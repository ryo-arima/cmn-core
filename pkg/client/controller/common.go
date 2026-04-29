package controller

import (
	"fmt"
	"os"
	"strings"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/spf13/cobra"
)

// global output format (table/json/yaml)
var outputFormat = "table"

// SetOutputFormat sets global output format
func SetOutputFormat(format string) {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "table", "json", "yaml":
		outputFormat = format
	default:
		outputFormat = "table"
	}
}

// GetOutputFormat returns current output format
func GetOutputFormat() string { return outputFormat }

// PrintMessage prints message as per current format via usecase formatter
func PrintMessage(msg string) {
	type message struct {
		Message string `json:"message" yaml:"message"`
	}
	fmt.Print(usecase.Format(GetOutputFormat(), message{Message: msg}))
}

// InitCommonRefreshTokenCmd creates a refresh command.
// Token refresh is not supported in this system; tokens must be re-issued by the IdP.
// The command informs the user and exits cleanly.
func InitCommonRefreshTokenCmd(manager *clientauth.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "refresh the access token",
		Long:  "token refresh is not supported; please re-authenticate via your IdP to obtain a new token",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.ForceRefresh(); err != nil {
				PrintMessage(err.Error())
			}
		},
	}
}

// InitCommonLogoutCmd creates a logout command that clears local token files.
func InitCommonLogoutCmd(manager *clientauth.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "logout and clear local tokens",
		Long:  "remove local token files; re-authentication via your IdP is required for future requests",
		Run: func(cmd *cobra.Command, args []string) {
			manager.ClearTokens()
			PrintMessage("local tokens cleared; please re-authenticate via your IdP")
		},
	}
}

// InitCommonValidateTokenCmd creates a validate command.
// The token is obtained transparently from the manager.
func InitCommonValidateTokenCmd(manager *clientauth.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "validate the current access token",
		Run: func(cmd *cobra.Command, args []string) {
			uc := usecase.NewCommon(manager.Conf(), manager)
			result := uc.ValidateToken()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// InitCommonUserInfoCmd creates a userinfo command.
// The token is obtained transparently from the manager.
func InitCommonUserInfoCmd(manager *clientauth.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "userinfo",
		Short: "get user information from the access token claims",
		Run: func(cmd *cobra.Command, args []string) {
			uc := usecase.NewCommon(manager.Conf(), manager)
			result := uc.GetUserInfo()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// InitSSOLoginCmd creates an explicit SSO login command.
func InitSSOLoginCmd(manager *clientauth.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "authenticate via SSO (OIDC)",
		Long:  "start a fresh SSO login flow via the configured OIDC provider.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.ForceLogin(""); err != nil {
				PrintMessage("login failed: " + err.Error())
				return
			}
			PrintMessage("login successful")
		},
	}
}
// InitAnonymousValidateCmd creates a validate command for the anonymous client.
// The token must be supplied via --access-token flag or CMN_ACCESS_TOKEN env var.
func InitAnonymousValidateCmd(conf config.BaseConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate an access token",
		Long:  "validate an externally-supplied access token against the server",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("access-token")
			if token == "" {
				token = os.Getenv("CMN_ACCESS_TOKEN")
			}
			if token == "" {
				PrintMessage("access-token is required (use --access-token or set CMN_ACCESS_TOKEN)")
				return
			}
			manager := clientauth.NewManager(conf, "anonymous").WithToken(token)
			uc := usecase.NewCommon(conf, manager)
			result := uc.ValidateToken()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
	cmd.Flags().StringP("access-token", "a", "", "access token to validate")
	return cmd
}

// NewDBBootstrapCmd returns a command that creates all required PostgreSQL
// tables via GORM AutoMigrate. Intended for first-time setup by an admin.
func NewDBBootstrapCmd(conf config.BaseConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "create all required database tables (PostgreSQL AutoMigrate)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := conf.ConnectDB(); err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			tables := []interface{}{
				&model.PgUsers{},
				&model.PgGroups{},
				&model.PgMembers{},
				&model.PgResource{},
				&model.PgResourceGroupRole{},
			}
			if err := conf.DBConnection.AutoMigrate(tables...); err != nil {
				return fmt.Errorf("AutoMigrate failed: %w", err)
			}
			fmt.Println("Bootstrap complete: all tables created or already up-to-date.")
			return nil
		},
	}
}
