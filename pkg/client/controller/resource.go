package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)

// InitResourceCmd returns the "resource" subcommand tree.
// It is shared by both the app client (/v1/internal) and admin client (/v1/private);
// the actual API prefix is determined by the usecase.ResourceUC implementation.
func InitResourceCmd(uc usecase.ResourceUC) *cobra.Command {
	resourceCmd := &cobra.Command{
		Use:   "resource",
		Short: "manage resources",
	}

	// list
	resourceCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list resources",
		Run: func(cmd *cobra.Command, args []string) {
			result := uc.ListResources()
			fmt.Print(usecase.Format(GetOutputFormat(), result))
		},
	})

	// get
	getCmd := &cobra.Command{
		Use:   "get",
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
	getCmd.Flags().String("uuid", "", "resource UUID")
	resourceCmd.AddCommand(getCmd)

	// create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create a new resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			desc, _ := cmd.Flags().GetString("description")
			result := uc.CreateResource(request.CreateResource{Name: name, Description: desc})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	createCmd.Flags().String("name", "", "resource name")
	createCmd.Flags().String("description", "", "resource description")
	resourceCmd.AddCommand(createCmd)

	// update
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			if uuid == "" {
				return fmt.Errorf("--uuid is required")
			}
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")
			result := uc.UpdateResource(uuid, request.UpdateResource{Name: name, Description: desc})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	updateCmd.Flags().String("uuid", "", "resource UUID")
	updateCmd.Flags().String("name", "", "new resource name")
	updateCmd.Flags().String("description", "", "new resource description")
	resourceCmd.AddCommand(updateCmd)

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete",
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
	deleteCmd.Flags().String("uuid", "", "resource UUID")
	resourceCmd.AddCommand(deleteCmd)

	// groups — subcommand group for resource group-role management
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "manage group-role assignments on a resource",
	}

	// groups list
	groupsListCmd := &cobra.Command{
		Use:   "list",
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
	groupsListCmd.Flags().String("uuid", "", "resource UUID")
	groupsCmd.AddCommand(groupsListCmd)

	// groups set
	groupsSetCmd := &cobra.Command{
		Use:   "set",
		Short: "assign a role to a group on a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			uuid, _ := cmd.Flags().GetString("uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")
			role, _ := cmd.Flags().GetString("role")
			if uuid == "" || groupUUID == "" || role == "" {
				return fmt.Errorf("--uuid, --group-uuid, and --role are required")
			}
			result := uc.SetResourceGroupRole(uuid, request.SetResourceGroupRole{
				GroupUUID: groupUUID,
				Role:      role,
			})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	groupsSetCmd.Flags().String("uuid", "", "resource UUID")
	groupsSetCmd.Flags().String("group-uuid", "", "group UUID")
	groupsSetCmd.Flags().String("role", "", "role: viewer, editor, or owner")
	groupsCmd.AddCommand(groupsSetCmd)

	// groups remove
	groupsRemoveCmd := &cobra.Command{
		Use:   "remove",
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
	groupsRemoveCmd.Flags().String("uuid", "", "resource UUID")
	groupsRemoveCmd.Flags().String("group-uuid", "", "group UUID")
	groupsCmd.AddCommand(groupsRemoveCmd)

	resourceCmd.AddCommand(groupsCmd)

	return resourceCmd
}
