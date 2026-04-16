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
