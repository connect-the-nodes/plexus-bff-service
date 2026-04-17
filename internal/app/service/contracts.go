package service

import (
	"context"
	"time"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/model"
)

type FeaturesRetriever interface {
	RetrieveFeatures(ctx context.Context) ([]model.FeatureFlag, error)
	IsActive(ctx context.Context, name string) (bool, error)
}

type ConnectorsService interface {
	GetConnectors(ctx context.Context) (dto.ConnectorsResponse, error)
}

type FeatureService interface {
	GetStatus() dto.FeatureResponse
}

type ObservabilityService interface {
	GetOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID string) (dto.ObservabilityOverviewResponse, error)
	GetInventory(ctx context.Context, from, to *time.Time, tenantID string) ([]dto.ObservabilityInventoryItemResponse, error)
	GetAdminOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID, connectorID string) (dto.AdminObservabilityOverviewResponse, error)
	GetAdminTenantHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.TenantHealthListResponse, error)
	GetAdminConnectorHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.ConnectorHealthListResponse, error)
}

type UserAdminService interface {
	ListUsers(ctx context.Context) ([]dto.PortalUser, error)
	GetUser(ctx context.Context, id string) (dto.PortalUser, error)
	CreateUser(ctx context.Context, input dto.CreatePortalUserInput) (dto.PortalUser, error)
	UpdateUser(ctx context.Context, id string, input dto.CreatePortalUserInput) (dto.PortalUser, error)
	DeleteUser(ctx context.Context, id string) error
}

type PermissionAdminService interface {
	ListPermissions(ctx context.Context) ([]dto.PermissionDefinition, error)
	GetPermission(ctx context.Context, id string) (dto.PermissionDefinition, error)
	CreatePermission(ctx context.Context, input dto.CreatePermissionInput) (dto.PermissionDefinition, error)
	UpdatePermission(ctx context.Context, id string, input dto.CreatePermissionInput) (dto.PermissionDefinition, error)
	DeletePermission(ctx context.Context, id string) error
}

type GroupAdminService interface {
	ListGroups(ctx context.Context) ([]dto.PortalGroup, error)
	GetGroup(ctx context.Context, id string) (dto.PortalGroup, error)
	CreateGroup(ctx context.Context, input dto.CreatePortalGroupInput) (dto.PortalGroup, error)
	UpdateGroup(ctx context.Context, id string, input dto.CreatePortalGroupInput) (dto.PortalGroup, error)
	DeleteGroup(ctx context.Context, id string) error
}

type DomainAdminService interface {
	GetWorkspace(ctx context.Context) (dto.DomainWorkspace, error)
	ListDomains(ctx context.Context) ([]dto.RegisteredDomain, error)
	ListApprovedDomains(ctx context.Context) ([]dto.RegisteredDomain, error)
	GetDomain(ctx context.Context, id string) (dto.RegisteredDomain, error)
	CreateDomain(ctx context.Context, input dto.DomainMutationRequest) (dto.RegisteredDomain, error)
	UpdateDomain(ctx context.Context, id string, input dto.DomainMutationRequest) (dto.RegisteredDomain, error)
	DeleteDomain(ctx context.Context, id string) error
	ReviewDomain(ctx context.Context, id string, input dto.ReviewDecisionRequest) (dto.RegisteredDomain, error)
}
