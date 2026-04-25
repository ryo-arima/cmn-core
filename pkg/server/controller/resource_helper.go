package controller

import (
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

func toResponseResource(r model.Resource) response.Resource {
	return response.Resource{
		ID:          r.ID,
		UUID:        r.UUID,
		Name:        r.Name,
		Description: r.Description,
		CreatedBy:   r.CreatedBy,
		UpdatedBy:   r.UpdatedBy,
		DeletedBy:   r.DeletedBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

func ptr[T any](v T) *T { return &v }
