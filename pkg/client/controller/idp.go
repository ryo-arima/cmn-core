package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)


// InitUserCmd returns the "user" subcommand tree for the app client.
// Commands call /v1/internal/user (own profile only).
func InitUserCmd(uc usecase.IdP) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "manage your own user profile",
	}

	// get
	getCmd := &cobra.Command{
		Use:   "get",
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
	getCmd.Flags().String("id", "", "user ID (omit for own profile)")
	userCmd.AddCommand(getCmd)

	// list
	userCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list users in your groups",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListGroupUsers()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	})

	// update
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update your own user profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var r request.UpdateUser
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
	updateCmd.Flags().String("email", "", "new email address")
	updateCmd.Flags().String("first-name", "", "new first name")
	updateCmd.Flags().String("last-name", "", "new last name")
	userCmd.AddCommand(updateCmd)

	return userCmd
}

// InitGroupCmd returns the "group" subcommand tree for the app client.
// Commands call /v1/internal/groups (groups the caller belongs to).
func InitGroupCmd(uc usecase.IdP) *cobra.Command {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "manage groups you belong to",
	}

	// list
	groupCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list groups you belong to",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListMyGroups()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	})

	// get
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get a group you belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			result := uc.GetGroup(id)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	getCmd.Flags().String("id", "", "group ID")
	groupCmd.AddCommand(getCmd)

	// create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create a new group",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			result := uc.CreateGroup(request.CreateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	createCmd.Flags().String("name", "", "group name")
	groupCmd.AddCommand(createCmd)

	// update
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update a group you belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			if id == "" || name == "" {
				return fmt.Errorf("--id and --name are required")
			}
			result := uc.UpdateGroup(id, request.UpdateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	updateCmd.Flags().String("id", "", "group ID")
	updateCmd.Flags().String("name", "", "new group name")
	groupCmd.AddCommand(updateCmd)

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a group you belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			result := uc.DeleteGroup(id)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	deleteCmd.Flags().String("id", "", "group ID")
	groupCmd.AddCommand(deleteCmd)

	return groupCmd
}

// InitMemberCmd returns the "member" subcommand tree for the app client.
func InitMemberCmd(uc usecase.IdP) *cobra.Command {
	memberCmd := &cobra.Command{
		Use:   "member",
		Short: "manage members of groups you belong to",
	}

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list members of a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("group-id")
			if gid == "" {
				return fmt.Errorf("--group-id is required")
			}
			result := uc.ListGroupMembers(gid)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	listCmd.Flags().String("group-id", "", "group ID")
	memberCmd.AddCommand(listCmd)

	// add
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add a user to a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("group-id")
			uid, _ := cmd.Flags().GetString("user-id")
			role, _ := cmd.Flags().GetString("role")
			if gid == "" || uid == "" {
				return fmt.Errorf("--group-id and --user-id are required")
			}
			if role == "" {
				return fmt.Errorf("--role is required (owner, editor, or viewer)")
			}
			result := uc.AddGroupMember(gid, request.AddGroupMember{UserID: uid, Role: role})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	addCmd.Flags().String("group-id", "", "group ID")
	addCmd.Flags().String("user-id", "", "user ID to add")
	addCmd.Flags().String("role", "", "role to assign: owner, editor, or viewer")
	memberCmd.AddCommand(addCmd)

	// remove
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "remove a user from a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("group-id")
			uid, _ := cmd.Flags().GetString("user-id")
			if gid == "" || uid == "" {
				return fmt.Errorf("--group-id and --user-id are required")
			}
			result := uc.RemoveGroupMember(gid, request.AddGroupMember{UserID: uid})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	removeCmd.Flags().String("group-id", "", "group ID")
	removeCmd.Flags().String("user-id", "", "user ID to remove")
	memberCmd.AddCommand(removeCmd)

	return memberCmd
}
