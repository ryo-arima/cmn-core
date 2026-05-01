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
	DeleteResourceGroupRole(uuid, groupID string) response.RrCommons
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

func (rcvr *resourceUsecase) ListResources() response.RrResources { return rcvr.repo.ListResources() }
func (rcvr *resourceUsecase) GetResource(uuid string) response.RrSingleResource {
	return rcvr.repo.GetResource(uuid)
}
func (rcvr *resourceUsecase) CreateResource(r request.RrCreateResource) response.RrSingleResource {
	return rcvr.repo.CreateResource(r)
}
func (rcvr *resourceUsecase) UpdateResource(uuid string, r request.RrUpdateResource) response.RrCommons {
	return rcvr.repo.UpdateResource(uuid, r)
}
func (rcvr *resourceUsecase) DeleteResource(uuid string) response.RrCommons {
	return rcvr.repo.DeleteResource(uuid)
}
func (rcvr *resourceUsecase) GetResourceGroupRoles(uuid string) response.RrResourceGroupRoles {
	return rcvr.repo.GetResourceGroupRoles(uuid)
}
func (rcvr *resourceUsecase) SetResourceGroupRole(uuid string, r request.RrSetResourceGroupRole) response.RrCommons {
	return rcvr.repo.SetResourceGroupRole(uuid, r)
}
func (rcvr *resourceUsecase) DeleteResourceGroupRole(uuid, groupID string) response.RrCommons {
	return rcvr.repo.DeleteResourceGroupRole(uuid, groupID)
}
