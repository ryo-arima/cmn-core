package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// Resource usecase interface.
type Resource interface {
	// List returns resources accessible to the caller (own + group).
	ListResources(ctx context.Context, userUUID string, groups []string) ([]model.PgResource, error)
	// Get returns a resource if the caller is allowed to view it.
	GetResource(ctx context.Context, resourceUUID, userUUID string, groups []string, isAdmin bool) (*model.PgResource, error)
	// Create creates a new resource owned by userID (IDP user ID), with an optional ownerGroup (IDP group ID).
	CreateResource(ctx context.Context, name, description, userID, ownerGroup string) (*model.PgResource, error)
	// Update modifies a resource if the caller has editor or owner role (or is admin/creator).
	UpdateResource(ctx context.Context, resourceUUID, name, description, userUUID string, groups []string, isAdmin bool) (*model.PgResource, error)
	// Delete soft-deletes a resource if the caller has owner role (or is admin/creator).
	DeleteResource(ctx context.Context, resourceUUID, userUUID string, groups []string, isAdmin bool) error

	// Admin-only: list all resources.
	ListAllResources(ctx context.Context) ([]model.PgResource, error)
	// Admin-only: hard delete.
	AdminDeleteResource(ctx context.Context, resourceUUID, userUUID string) error

	// Group role management (owner or creator of the resource required).
	GetGroupRoles(ctx context.Context, resourceUUID, userID string, groups []string, isAdmin bool) ([]model.PgResourceGroupRole, error)
	SetGroupRole(ctx context.Context, resourceUUID, groupID, role, userID string, groups []string, isAdmin bool) error
	DeleteGroupRole(ctx context.Context, resourceUUID, groupID, userID string, groups []string, isAdmin bool) error
}

type resourceUsecase struct {
	repo repository.Resource
}

// NewResource creates a new Resource usecase.
func NewResource(repo repository.Resource) Resource {
	return &resourceUsecase{repo: repo}
}

// ---- helpers ---------------------------------------------------------------

// canView returns true if the caller may view the resource.
func canView(res *model.PgResource, userUUID string, groups []string, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	if res.CreatedBy == userUUID {
		return true
	}
	return hasGroupRole(res.UUID, groups) // group lookup is done at DB level; we check non-empty groups here
}

// hasGroupRoleInList reports whether any of the caller's groups appears in the provided role list.
func hasGroupRoleInList(roles []model.PgResourceGroupRole, groups []string, minRole string) bool {
	allowed := roleLevel(minRole)
	for _, r := range roles {
		if roleLevel(r.Role) >= allowed {
			for _, g := range groups {
				if r.GroupID == g {
					return true
				}
			}
		}
	}
	return false
}

// hasGroupRole is a simple membership check (any role grants view permission).
func hasGroupRole(_ string, _ []string) bool {
	// The DB query already filters; this is a no-op guard used in canView for non-DB paths.
	return false
}

func roleLevel(role string) int {
	switch role {
	case "owner":
		return 3
	case "editor":
		return 2
	case "viewer":
		return 1
	}
	return 0
}

// canManage checks if the caller can modify or delete the resource (editor+) or manage groups (owner only).
func canManage(roles []model.PgResourceGroupRole, res *model.PgResource, userUUID string, groups []string, isAdmin bool, minRole string) bool {
	if isAdmin {
		return true
	}
	if res.CreatedBy == userUUID {
		return true
	}
	return hasGroupRoleInList(roles, groups, minRole)
}

// ---- implementation --------------------------------------------------------

func (uc *resourceUsecase) ListResources(ctx context.Context, userUUID string, groups []string) ([]model.PgResource, error) {
	return uc.repo.ListResources(ctx, request.LoResourceQueryFilter{
		CreatedBy: userUUID,
		GroupIDs:  groups,
	})
}

func (uc *resourceUsecase) GetResource(ctx context.Context, resourceUUID, userUUID string, groups []string, isAdmin bool) (*model.PgResource, error) {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	if isAdmin || res.CreatedBy == userUUID {
		return res, nil
	}
	// Check group membership via DB
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	if hasGroupRoleInList(roles, groups, "viewer") {
		return res, nil
	}
	return nil, fmt.Errorf("access denied")
}

func (uc *resourceUsecase) CreateResource(ctx context.Context, name, description, userID, ownerGroup string) (*model.PgResource, error) {
	res := &model.PgResource{
		UUID:        uuid.New().String(),
		Name:        name,
		Description: description,
		OwnerGroup:  ownerGroup,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	if err := uc.repo.CreateResource(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *resourceUsecase) UpdateResource(ctx context.Context, resourceUUID, name, description, userUUID string, groups []string, isAdmin bool) (*model.PgResource, error) {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	if !canManage(roles, res, userUUID, groups, isAdmin, "editor") {
		return nil, fmt.Errorf("access denied")
	}
	if name != "" {
		res.Name = name
	}
	if description != "" {
		res.Description = description
	}
	res.UpdatedBy = userUUID
	if err := uc.repo.UpdateResource(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *resourceUsecase) DeleteResource(ctx context.Context, resourceUUID, userUUID string, groups []string, isAdmin bool) error {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return err
	}
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return err
	}
	if !canManage(roles, res, userUUID, groups, isAdmin, "owner") {
		return fmt.Errorf("access denied")
	}
	return uc.repo.SoftDeleteResource(ctx, res, userUUID)
}

func (uc *resourceUsecase) ListAllResources(ctx context.Context) ([]model.PgResource, error) {
	return uc.repo.ListAllResources(ctx)
}

func (uc *resourceUsecase) AdminDeleteResource(ctx context.Context, resourceUUID, userUUID string) error {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return err
	}
	return uc.repo.SoftDeleteResource(ctx, res, userUUID)
}

func (uc *resourceUsecase) GetGroupRoles(ctx context.Context, resourceUUID, userID string, groups []string, isAdmin bool) ([]model.PgResourceGroupRole, error) {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return nil, err
	}
	if !isAdmin && res.CreatedBy != userID && !hasGroupRoleInList(roles, groups, "viewer") {
		return nil, fmt.Errorf("access denied")
	}
	return roles, nil
}

func (uc *resourceUsecase) SetGroupRole(ctx context.Context, resourceUUID, groupID, role, userID string, groups []string, isAdmin bool) error {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return err
	}
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return err
	}
	if !canManage(roles, res, userID, groups, isAdmin, "owner") {
		return fmt.Errorf("access denied")
	}
	return uc.repo.SetGroupRole(ctx, &model.PgResourceGroupRole{
		ResourceUUID: resourceUUID,
		GroupID:      groupID,
		Role:         role,
	})
}

func (uc *resourceUsecase) DeleteGroupRole(ctx context.Context, resourceUUID, groupID, userID string, groups []string, isAdmin bool) error {
	res, err := uc.repo.GetResourceByUUID(ctx, resourceUUID)
	if err != nil {
		return err
	}
	roles, err := uc.repo.GetGroupRoles(ctx, resourceUUID)
	if err != nil {
		return err
	}
	if !canManage(roles, res, userID, groups, isAdmin, "owner") {
		return fmt.Errorf("access denied")
	}
	return uc.repo.DeleteGroupRole(ctx, resourceUUID, groupID)
}
