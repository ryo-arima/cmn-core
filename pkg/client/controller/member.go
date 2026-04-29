package controller

import (
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/client/usecase"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/spf13/cobra"
)

// ---- Member commands (app client: /v1/internal) ----

// NewMemberListCmd returns the "member" subcommand for use under "list".
func NewMemberListCmd(uc usecase.Member) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
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
	cmd.Flags().String("group-id", "", "group ID")
	return cmd
}

// NewMemberAddCmd returns the "member" subcommand for use under "add".
func NewMemberAddCmd(uc usecase.Member) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
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
			result := uc.AddGroupMember(gid, request.RrAddGroupMember{UserID: uid, Role: role})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("group-id", "", "group ID")
	cmd.Flags().String("user-id", "", "user ID to add")
	cmd.Flags().String("role", "", "role to assign: owner, editor, or viewer")
	return cmd
}

// NewMemberRemoveCmd returns the "member" subcommand for use under "remove".
func NewMemberRemoveCmd(uc usecase.Member) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "remove a user from a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("group-id")
			uid, _ := cmd.Flags().GetString("user-id")
			if gid == "" || uid == "" {
				return fmt.Errorf("--group-id and --user-id are required")
			}
			result := uc.RemoveGroupMember(gid, request.RrRemoveGroupMember{UserID: uid})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("group-id", "", "group ID")
	cmd.Flags().String("user-id", "", "user ID to remove")
	return cmd
}

// ---- Admin member commands (admin client: /v1/private) ----

// NewAdminMemberListCmd returns the "member" subcommand for use under "list" (admin).
func NewAdminMemberListCmd(uc usecase.MemberAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
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
	cmd.Flags().String("group-id", "", "group ID")
	return cmd
}

// NewAdminMemberAddCmd returns the "member" subcommand for use under "add" (admin).
func NewAdminMemberAddCmd(uc usecase.MemberAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
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
			result := uc.AddGroupMember(gid, request.RrAddGroupMember{UserID: uid, Role: role})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("group-id", "", "group ID")
	cmd.Flags().String("user-id", "", "user ID to add")
	cmd.Flags().String("role", "", "role to assign: owner, editor, or viewer")
	return cmd
}

// NewAdminMemberRemoveCmd returns the "member" subcommand for use under "remove" (admin).
func NewAdminMemberRemoveCmd(uc usecase.MemberAdmin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "remove a user from a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("group-id")
			uid, _ := cmd.Flags().GetString("user-id")
			if gid == "" || uid == "" {
				return fmt.Errorf("--group-id and --user-id are required")
			}
			result := uc.RemoveGroupMember(gid, request.RrRemoveGroupMember{UserID: uid})
			fmt.Print(usecase.Format(GetOutputFormat(), result))
			return nil
		},
	}
	cmd.Flags().String("group-id", "", "group ID")
	cmd.Flags().String("user-id", "", "user ID to remove")
	return cmd
}
