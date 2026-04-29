package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// ResourceUC is the business-logic interface for resource operations.
type ResourceUC interface {
	ListResources() response.RrResources
	GetResource(uuid string) response.RrSingleResource
	CreateResource(req request.RrCreateResource) response.RrSingleResource
	UpdateResource(uuid string, req request.RrUpdateResource) response.RrCommons
	DeleteResource(uuid string) response.RrCommons
	GetResourceGroupRoles(uuid string) response.RrResourceGroupRoles
	SetResourceGroupRole(uuid string, req request.RrSetResourceGroupRole) response.RrCommons
	DeleteResourceGroupRole(uuid, groupUUID string) response.RrCommons
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

func (u *resourceUsecase) ListResources() response.RrResources { return u.repo.ListResources() }
func (u *resourceUsecase) GetResource(uuid string) response.RrSingleResource {
	return u.repo.GetResource(uuid)
}
func (u *resourceUsecase) CreateResource(r request.RrCreateResource) response.RrSingleResource {
	return u.repo.CreateResource(r)
}
func (u *resourceUsecase) UpdateResource(uuid string, r request.RrUpdateResource) response.RrCommons {
	return u.repo.UpdateResource(uuid, r)
}
func (u *resourceUsecase) DeleteResource(uuid string) response.RrCommons {
	return u.repo.DeleteResource(uuid)
}
func (u *resourceUsecase) GetResourceGroupRoles(uuid string) response.RrResourceGroupRoles {
	return u.repo.GetResourceGroupRoles(uuid)
}
func (u *resourceUsecase) SetResourceGroupRole(uuid string, r request.RrSetResourceGroupRole) response.RrCommons {
	return u.repo.SetResourceGroupRole(uuid, r)
}
func (u *resourceUsecase) DeleteResourceGroupRole(uuid, groupUUID string) response.RrCommons {
	return u.repo.DeleteResourceGroupRole(uuid, groupUUID)
}
