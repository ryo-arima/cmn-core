package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"gorm.io/gorm"
)

type resourceRepository struct {
	db *gorm.DB
}

// NewResource creates a new Resource repository backed by the given database connection.
func NewResource(conf config.BaseConfig) Resource {
	return &resourceRepository{db: conf.DBConnection}
}

func (r *resourceRepository) GetResourceByUUID(ctx context.Context, uuid string) (*model.Resource, error) {
	var res model.Resource
	if err := r.db.WithContext(ctx).Where("uuid = ? AND deleted_at IS NULL", uuid).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("resource not found: %s", uuid)
		}
		return nil, err
	}
	return &res, nil
}

// ListResources returns resources accessible to the caller:
//   - all resources whose created_by matches filter.CreatedBy, or
//   - all resources that have a ResourceGroupRole entry for any of filter.GroupUUIDs.
func (r *resourceRepository) ListResources(ctx context.Context, filter ResourceQueryFilter) ([]model.Resource, error) {
	var resources []model.Resource

	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if len(filter.GroupUUIDs) > 0 {
		// created_by matches OR resource has a group-role entry for one of the user's groups
		query = query.Where(
			"created_by = ? OR uuid IN (SELECT resource_uuid FROM resource_group_roles WHERE group_uuid IN ?)",
			filter.CreatedBy, filter.GroupUUIDs,
		)
	} else {
		query = query.Where("created_by = ?", filter.CreatedBy)
	}

	if err := query.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// ListAllResources returns every non-deleted resource (admin only).
func (r *resourceRepository) ListAllResources(ctx context.Context) ([]model.Resource, error) {
	var resources []model.Resource
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *resourceRepository) CreateResource(ctx context.Context, resource *model.Resource) error {
	return r.db.WithContext(ctx).Create(resource).Error
}

func (r *resourceRepository) UpdateResource(ctx context.Context, resource *model.Resource) error {
	return r.db.WithContext(ctx).Save(resource).Error
}

func (r *resourceRepository) SoftDeleteResource(ctx context.Context, resource *model.Resource, deletedBy string) error {
	now := time.Now()
	resource.DeletedBy = deletedBy
	resource.DeletedAt = &now
	return r.db.WithContext(ctx).Save(resource).Error
}

func (r *resourceRepository) GetGroupRoles(ctx context.Context, resourceUUID string) ([]model.ResourceGroupRole, error) {
	var roles []model.ResourceGroupRole
	if err := r.db.WithContext(ctx).Where("resource_uuid = ?", resourceUUID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// SetGroupRole inserts or updates (upsert) a group-role entry for a resource.
func (r *resourceRepository) SetGroupRole(ctx context.Context, rgr *model.ResourceGroupRole) error {
	var existing model.ResourceGroupRole
	err := r.db.WithContext(ctx).
		Where("resource_uuid = ? AND group_uuid = ?", rgr.ResourceUUID, rgr.GroupUUID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.WithContext(ctx).Create(rgr).Error
	}
	if err != nil {
		return err
	}
	existing.Role = rgr.Role
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *resourceRepository) DeleteGroupRole(ctx context.Context, resourceUUID, groupUUID string) error {
	return r.db.WithContext(ctx).
		Where("resource_uuid = ? AND group_uuid = ?", resourceUUID, groupUUID).
		Delete(&model.ResourceGroupRole{}).Error
}
