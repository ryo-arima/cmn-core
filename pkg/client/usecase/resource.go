package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// ResourceUC is the business-logic interface for resource operations.
type ResourceUC interface {
	ListResources() response.Resources
	GetResource(uuid string) response.SingleResource
	CreateResource(req request.CreateResource) response.SingleResource
	UpdateResource(uuid string, req request.UpdateResource) response.Commons
	DeleteResource(uuid string) response.Commons
	GetResourceGroupRoles(uuid string) response.ResourceGroupRoles
	SetResourceGroupRole(uuid string, req request.SetResourceGroupRole) response.Commons
	DeleteResourceGroupRole(uuid, groupUUID string) response.Commons
}

type resourceUsecase struct {
	repo repository.Resource
}

// NewResourceUC creates a ResourceUC usecase for the app client (/v1/internal).
func NewResourceUC(conf config.BaseConfig, manager *clientauth.Manager) ResourceUC {
	return &resourceUsecase{repo: repository.NewResource(conf, manager)}
}

// NewResourceAdminUC creates a ResourceUC usecase for the admin client (/v1/private).
func NewResourceAdminUC(conf config.BaseConfig, manager *clientauth.Manager) ResourceUC {
	return &resourceUsecase{repo: repository.NewResourceAdmin(conf, manager)}
}

func (u *resourceUsecase) ListResources() response.Resources { return u.repo.ListResources() }
func (u *resourceUsecase) GetResource(uuid string) response.SingleResource {
	return u.repo.GetResource(uuid)
}
func (u *resourceUsecase) CreateResource(r request.CreateResource) response.SingleResource {
	return u.repo.CreateResource(r)
}
func (u *resourceUsecase) UpdateResource(uuid string, r request.UpdateResource) response.Commons {
	return u.repo.UpdateResource(uuid, r)
}
func (u *resourceUsecase) DeleteResource(uuid string) response.Commons {
	return u.repo.DeleteResource(uuid)
}
func (u *resourceUsecase) GetResourceGroupRoles(uuid string) response.ResourceGroupRoles {
	return u.repo.GetResourceGroupRoles(uuid)
}
func (u *resourceUsecase) SetResourceGroupRole(uuid string, r request.SetResourceGroupRole) response.Commons {
	return u.repo.SetResourceGroupRole(uuid, r)
}
func (u *resourceUsecase) DeleteResourceGroupRole(uuid, groupUUID string) response.Commons {
	return u.repo.DeleteResourceGroupRole(uuid, groupUUID)
}
