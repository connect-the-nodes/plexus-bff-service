package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/security"
)

type fakeObservabilityService struct {
	lastTenant string
}

func (f *fakeObservabilityService) GetOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID string) (dto.ObservabilityOverviewResponse, error) {
	f.lastTenant = tenantID
	return dto.ObservabilityOverviewResponse{}, nil
}

func (f *fakeObservabilityService) GetInventory(ctx context.Context, from, to *time.Time, tenantID string) ([]dto.ObservabilityInventoryItemResponse, error) {
	f.lastTenant = tenantID
	return []dto.ObservabilityInventoryItemResponse{}, nil
}

func (f *fakeObservabilityService) GetAdminOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID, connectorID string) (dto.AdminObservabilityOverviewResponse, error) {
	return dto.AdminObservabilityOverviewResponse{}, nil
}

func (f *fakeObservabilityService) GetAdminTenantHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.TenantHealthListResponse, error) {
	return dto.TenantHealthListResponse{}, nil
}

func (f *fakeObservabilityService) GetAdminConnectorHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.ConnectorHealthListResponse, error) {
	return dto.ConnectorHealthListResponse{}, nil
}

func TestObservabilityUsesTenantFromJWTClaims(t *testing.T) {
	service := &fakeObservabilityService{}
	controller := NewObservabilityController(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v3/observability/overview", nil)
	request = request.WithContext(security.WithAuthentication(request.Context(), &security.Authentication{
		Claims: map[string]any{"tenantId": "tenant-jwt"},
	}))
	recorder := httptest.NewRecorder()

	controller.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if service.lastTenant != "tenant-jwt" {
		t.Fatalf("expected tenant tenant-jwt, got %s", service.lastTenant)
	}
}

func TestObservabilityRejectsTenantMismatch(t *testing.T) {
	service := &fakeObservabilityService{}
	controller := NewObservabilityController(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v3/observability/overview?tenantId=tenant-query", nil)
	request = request.WithContext(security.WithAuthentication(request.Context(), &security.Authentication{
		Claims: map[string]any{"tenantId": "tenant-jwt"},
	}))
	recorder := httptest.NewRecorder()

	controller.Overview(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestObservabilityFallsBackToDefaultTenant(t *testing.T) {
	service := &fakeObservabilityService{}
	controller := NewObservabilityController(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v3/observability/integrations/inventory", nil)
	recorder := httptest.NewRecorder()

	controller.Inventory(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if service.lastTenant != "default" {
		t.Fatalf("expected default tenant, got %s", service.lastTenant)
	}
}
