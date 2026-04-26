package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)
// InitAdminUserCmd returns the "user" subcommand tree for the admin client.
// Commands call /v1/private/users (all users).
func InitAdminUserCmd(uc usecase.IdPAdmin) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "manage users via the admin API",
	}

	// list
	userCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list all users",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListUsers()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	})

	// get
	getCmd := &cobra.Command{
		Use:   "get",
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
	getCmd.Flags().String("id", "", "user ID")
	userCmd.AddCommand(getCmd)

	// create
	createCmd := &cobra.Command{
		Use:   "create",
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
			result := uc.CreateUser(request.CreateUser{
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
	createCmd.Flags().String("username", "", "username")
	createCmd.Flags().String("email", "", "email address")
	createCmd.Flags().String("password", "", "initial password")
	createCmd.Flags().String("first-name", "", "first name")
	createCmd.Flags().String("last-name", "", "last name")
	userCmd.AddCommand(createCmd)

	// update
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required")
			}
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
			result := uc.UpdateUser(id, r)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	updateCmd.Flags().String("id", "", "user ID")
	updateCmd.Flags().String("email", "", "new email address")
	updateCmd.Flags().String("first-name", "", "new first name")
	updateCmd.Flags().String("last-name", "", "new last name")
	userCmd.AddCommand(updateCmd)

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete",
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
	deleteCmd.Flags().String("id", "", "user ID")
	userCmd.AddCommand(deleteCmd)

	return userCmd
}

// InitAdminGroupCmd returns the "group" subcommand tree for the admin client.
// Commands call /v1/private/groups (all groups).
func InitAdminGroupCmd(uc usecase.IdPAdmin) *cobra.Command {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "manage groups via the admin API",
	}

	// list
	groupCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list all groups",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListGroups()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	})

	// get
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get a group by ID",
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
		Short: "update a group",
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
		Short: "delete a group",
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

// InitAdminMemberCmd returns the "member" subcommand tree for the admin client.
func InitAdminMemberCmd(uc usecase.IdPAdmin) *cobra.Command {
	memberCmd := &cobra.Command{
		Use:   "member",
		Short: "manage group memberships via the admin API",
	}

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
