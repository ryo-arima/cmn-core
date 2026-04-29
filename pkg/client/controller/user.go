package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)

// ---- User commands (app client: /v1/internal) ----

// NewUserGetCmd returns the "user" subcommand for use under "get".
func NewUserGetCmd(uc usecase.User) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "get your own profile, or any user by --id",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			if id != "" {
				result := uc.GetUser(id)
				if result.User != nil {
					u := result.User
					fmt.Printf("%-20s %s\n", "ID:", u.ID)
					fmt.Printf("%-20s %s\n", "Username:", u.Username)
					fmt.Printf("%-20s %s\n", "Email:", u.Email)
					fmt.Printf("%-20s %s\n", "First name:", u.FirstName)
					fmt.Printf("%-20s %s\n", "Last name:", u.LastName)
					fmt.Printf("%-20s %v\n", "Enabled:", u.Enabled)
				} else {
					fmt.Printf("[%s] %s\n", result.Code, result.Message)
				}
			} else {
				result := uc.GetMyUser()
				if result.User != nil {
					u := result.User
					fmt.Printf("%-20s %s\n", "ID:", u.ID)
					fmt.Printf("%-20s %s\n", "Username:", u.Username)
					fmt.Printf("%-20s %s\n", "Email:", u.Email)
					fmt.Printf("%-20s %s\n", "First name:", u.FirstName)
					fmt.Printf("%-20s %s\n", "Last name:", u.LastName)
					fmt.Printf("%-20s %v\n", "Enabled:", u.Enabled)
				} else {
					fmt.Printf("[%s] %s\n", result.Code, result.Message)
				}
			}
		},
	}
	cmd.Flags().String("id", "", "user ID (omit for own profile)")
	return cmd
}

// NewUserListCmd returns the "user" subcommand for use under "list".
func NewUserListCmd(uc usecase.User) *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "list users in your groups",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListGroupUsers()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// NewUserUpdateCmd returns the "user" subcommand for use under "update".
func NewUserUpdateCmd(uc usecase.User) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "update your own user profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var r request.RrUpdateUser
			if email, _ := cmd.Flags().GetString("email"); email != "" {
				r.Email = &email
			}
			if fn, _ := cmd.Flags().GetString("first-name"); fn != "" {
				r.FirstName = &fn
			}
			if ln, _ := cmd.Flags().GetString("last-name"); ln != "" {
				r.LastName = &ln
			}
			result := uc.UpdateMyUser(r)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("email", "", "new email address")
	cmd.Flags().String("first-name", "", "new first name")
	cmd.Flags().String("last-name", "", "new last name")
	return cmd
}

// ---- Admin user commands (admin client: /v1/private/users) ----

// NewAdminUserGetCmd returns the "user" subcommand for use under "get" (admin).
func NewAdminUserGetCmd(uc usecase.UserAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "get a user by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			result := uc.GetUser(id)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("id", "", "user ID")
	return cmd
}

// NewAdminUserListCmd returns the "user" subcommand for use under "list" (admin).
func NewAdminUserListCmd(uc usecase.UserAdmin) *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "list all users",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListUsers()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// NewAdminUserCreateCmd returns the "user" subcommand for use under "create" (admin).
func NewAdminUserCreateCmd(uc usecase.UserAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")
			if username == "" || email == "" || password == "" {
				return fmt.Errorf("--username, --email, and --password are required")
			}
			fn, _ := cmd.Flags().GetString("first-name")
			ln, _ := cmd.Flags().GetString("last-name")
			result := uc.CreateUser(request.RrCreateUser{
				Username:  username,
				Email:     email,
				Password:  password,
				FirstName: fn,
				LastName:  ln,
				Enabled:   true,
			})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("username", "", "username")
	cmd.Flags().String("email", "", "email address")
	cmd.Flags().String("password", "", "initial password")
	cmd.Flags().String("first-name", "", "first name")
	cmd.Flags().String("last-name", "", "last name")
	return cmd
}

// NewAdminUserUpdateCmd returns the "user" subcommand for use under "update" (admin).
func NewAdminUserUpdateCmd(uc usecase.UserAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "update a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			var r request.RrUpdateUser
			if email, _ := cmd.Flags().GetString("email"); email != "" {
				r.Email = &email
			}
			if fn, _ := cmd.Flags().GetString("first-name"); fn != "" {
				r.FirstName = &fn
			}
			if ln, _ := cmd.Flags().GetString("last-name"); ln != "" {
				r.LastName = &ln
			}
			result := uc.UpdateUser(id, r)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("id", "", "user ID")
	cmd.Flags().String("email", "", "new email address")
	cmd.Flags().String("first-name", "", "new first name")
	cmd.Flags().String("last-name", "", "new last name")
	return cmd
}

// NewAdminUserDeleteCmd returns the "user" subcommand for use under "delete" (admin).
func NewAdminUserDeleteCmd(uc usecase.UserAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "delete a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			result := uc.DeleteUser(id)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("id", "", "user ID")
	return cmd
}

// ---- Anonymous user commands (/v1/public) ----

// NewAnonymousCreateUserCmd returns the "user" subcommand for use under "create"
// in the anonymous client. Calls POST /v1/public/user (no authentication required).
func NewAnonymousCreateUserCmd(uc usecase.Anonymous) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "register a new user (public endpoint, no authentication required)",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")
			if username == "" || email == "" || password == "" {
				return fmt.Errorf("--username, --email, and --password are required")
			}
			fn, _ := cmd.Flags().GetString("first-name")
			ln, _ := cmd.Flags().GetString("last-name")
			result := uc.RegisterUser(request.RrCreateUser{
				Username:  username,
				Email:     email,
				Password:  password,
				FirstName: fn,
				LastName:  ln,
				Enabled:   true,
			})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("username", "", "username")
	cmd.Flags().String("email", "", "email address")
	cmd.Flags().String("password", "", "initial password")
	cmd.Flags().String("first-name", "", "first name (optional)")
	cmd.Flags().String("last-name", "", "last name (optional)")
	return cmd
}
