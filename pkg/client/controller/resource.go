package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)

// ---- Resource commands ----
// Shared by the app client (/v1/internal) and admin client (/v1/private);
// the actual API prefix is determined by the usecase.ResourceUC implementation.

// NewResourceGetCmd returns the "resource" subcommand for use under "get".
func NewResourceGetCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "get a resource by UUID",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			if uuid == "" {
				return fmt.Errorf("--uuid is required")
			}
			result := uc.GetResource(uuid)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	return cmd
}

// NewResourceListCmd returns the "resource" subcommand for use under "list".
func NewResourceListCmd(uc usecase.ResourceUC) *cobra.Command {
	return &cobra.Command{
		Use:   "resource",
		Short: "list resources",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListResources()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	}
}

// NewResourceCreateCmd returns the "resource" subcommand for use under "create".
func NewResourceCreateCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "create a new resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			desc, _ := cmd.Flags().GetString("description")
			result := uc.CreateResource(request.RrCreateResource{Name: name, Description: desc})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("name", "", "resource name")
	cmd.Flags().String("description", "", "resource description")
	return cmd
}

// NewResourceUpdateCmd returns the "resource" subcommand for use under "update".
func NewResourceUpdateCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "update a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			if uuid == "" {
				return fmt.Errorf("--uuid is required")
			}
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")
			result := uc.UpdateResource(uuid, request.RrUpdateResource{Name: name, Description: desc})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	cmd.Flags().String("name", "", "new resource name")
	cmd.Flags().String("description", "", "new resource description")
	return cmd
}

// NewResourceDeleteCmd returns the "resource" subcommand for use under "delete".
func NewResourceDeleteCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "delete a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			if uuid == "" {
				return fmt.Errorf("--uuid is required")
			}
			result := uc.DeleteResource(uuid)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	return cmd
}

// ---- Resource-group commands ----

// NewResourceGroupListCmd returns the "resource-group" subcommand for use under "list".
func NewResourceGroupListCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-group",
		Short: "list group-role assignments for a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			if uuid == "" {
				return fmt.Errorf("--uuid is required")
			}
			result := uc.GetResourceGroupRoles(uuid)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	return cmd
}

// NewResourceGroupSetCmd returns the "resource-group" subcommand for use under "set".
func NewResourceGroupSetCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-group",
		Short: "assign a role to a group on a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")
			role, _ := cmd.Flags().GetString("role")
			if uuid == "" || groupUUID == "" || role == "" {
				return fmt.Errorf("--uuid, --group-uuid, and --role are required")
			}
			result := uc.SetResourceGroupRole(uuid, request.RrSetResourceGroupRole{
				GroupUUID: groupUUID,
				Role:      role,
			})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	cmd.Flags().String("group-uuid", "", "group UUID")
	cmd.Flags().String("role", "", "role: viewer, editor, or owner")
	return cmd
}

// NewResourceGroupRemoveCmd returns the "resource-group" subcommand for use under "remove".
func NewResourceGroupRemoveCmd(uc usecase.ResourceUC) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-group",
		Short: "remove a group from a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")
			if uuid == "" || groupUUID == "" {
				return fmt.Errorf("--uuid and --group-uuid are required")
			}
			result := uc.DeleteResourceGroupRole(uuid, groupUUID)
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("uuid", "", "resource UUID")
	cmd.Flags().String("group-uuid", "", "group UUID")
	return cmd
}
