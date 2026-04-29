package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"gorm.io/gorm"
)

// Group is a repository for persisting group UUID ↔ display name mappings.
// Used when the IdP (e.g. Casdoor) does not assign native UUIDs to groups.
type Group interface {
	Upsert(ctx context.Context, uuid, name string) error
	LookupName(ctx context.Context, uuid string) string
	LookupNames(ctx context.Context, uuids []string) map[string]string
	SoftDelete(ctx context.Context, uuid string) error
}

// NewGroup creates a PostgreSQL-backed Group repository.
func NewGroup(db *gorm.DB) Group {
	return &groupRepository{db: db}
}

type groupRepository struct {
	db *gorm.DB
}

func (r *groupRepository) Upsert(ctx context.Context, uuid, name string) error {
	now := time.Now()
	var row model.PgGroups
	err := r.db.WithContext(ctx).Where("uuid = ? AND deleted_at IS NULL", uuid).First(&row).Error
	if err == nil {
		return r.db.WithContext(ctx).Model(&row).Updates(map[string]interface{}{
			"name":       name,
			"updated_at": &now,
		}).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = model.PgGroups{UUID: uuid, Name: name, CreatedAt: &now, UpdatedAt: &now}
		return r.db.WithContext(ctx).Create(&row).Error
	}
	return err
}

func (r *groupRepository) LookupName(ctx context.Context, uuid string) string {
	var row model.PgGroups
	if err := r.db.WithContext(ctx).Where("uuid = ? AND deleted_at IS NULL", uuid).First(&row).Error; err != nil {
		return uuid
	}
	return row.Name
}

func (r *groupRepository) LookupNames(ctx context.Context, uuids []string) map[string]string {
	result := make(map[string]string, len(uuids))
	if len(uuids) == 0 {
		return result
	}
	var rows []model.PgGroups
	if err := r.db.WithContext(ctx).Where("uuid IN ? AND deleted_at IS NULL", uuids).Find(&rows).Error; err == nil {
		for _, row := range rows {
			result[row.UUID] = row.Name
		}
	}
	for _, u := range uuids {
		if _, ok := result[u]; !ok {
			result[u] = u
		}
	}
	return result
}

func (r *groupRepository) SoftDelete(ctx context.Context, uuid string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.PgGroups{}).
		Where("uuid = ? AND deleted_at IS NULL", uuid).
		Update("deleted_at", &now).Error
}
