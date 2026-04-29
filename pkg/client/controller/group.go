package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)

// ---- Group commands (app client: /v1/internal/groups) ----

// NewGroupGetCmd returns the "group" subcommand for use under "get".
func NewGroupGetCmd(uc usecase.Group) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
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
	cmd.Flags().String("id", "", "group ID")
	return cmd
}

// NewGroupListCmd returns the "group" subcommand for use under "list".
func NewGroupListCmd(uc usecase.Group) *cobra.Command {
	return &cobra.Command{
		Use:   "group",
		Short: "list groups you belong to",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListMyGroups()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// NewGroupCreateCmd returns the "group" subcommand for use under "create".
func NewGroupCreateCmd(uc usecase.Group) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "create a new group",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			result := uc.CreateGroup(request.RrCreateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("name", "", "group name")
	return cmd
}

// NewGroupUpdateCmd returns the "group" subcommand for use under "update".
func NewGroupUpdateCmd(uc usecase.Group) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "update a group you belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			if id == "" || name == "" {
				return fmt.Errorf("--id and --name are required")
			}
			result := uc.UpdateGroup(id, request.RrUpdateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("id", "", "group ID")
	cmd.Flags().String("name", "", "new group name")
	return cmd
}

// NewGroupDeleteCmd returns the "group" subcommand for use under "delete".
func NewGroupDeleteCmd(uc usecase.Group) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
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
	cmd.Flags().String("id", "", "group ID")
	return cmd
}

// ---- Admin group commands (admin client: /v1/private/groups) ----

// NewAdminGroupGetCmd returns the "group" subcommand for use under "get" (admin).
func NewAdminGroupGetCmd(uc usecase.GroupAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
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
	cmd.Flags().String("id", "", "group ID")
	return cmd
}

// NewAdminGroupListCmd returns the "group" subcommand for use under "list" (admin).
func NewAdminGroupListCmd(uc usecase.GroupAdmin) *cobra.Command {
	return &cobra.Command{
		Use:   "group",
		Short: "list all groups",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListGroups()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// NewAdminGroupCreateCmd returns the "group" subcommand for use under "create" (admin).
func NewAdminGroupCreateCmd(uc usecase.GroupAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "create a new group",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			result := uc.CreateGroup(request.RrCreateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("name", "", "group name")
	return cmd
}

// NewAdminGroupUpdateCmd returns the "group" subcommand for use under "update" (admin).
func NewAdminGroupUpdateCmd(uc usecase.GroupAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "update a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			if id == "" || name == "" {
				return fmt.Errorf("--id and --name are required")
			}
			result := uc.UpdateGroup(id, request.RrUpdateGroup{Name: name})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("id", "", "group ID")
	cmd.Flags().String("name", "", "new group name")
	return cmd
}

// NewAdminGroupDeleteCmd returns the "group" subcommand for use under "delete" (admin).
func NewAdminGroupDeleteCmd(uc usecase.GroupAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
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
	cmd.Flags().String("id", "", "group ID")
	return cmd
}
