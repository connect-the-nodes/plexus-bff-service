package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type ObservabilityService struct {
	baseURL string
	client  *http.Client
}

func NewObservabilityService(baseURL string) service.ObservabilityService {
	return &ObservabilityService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ObservabilityService) GetOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID string) (dto.ObservabilityOverviewResponse, error) {
	var result dto.ObservabilityOverviewResponse
	err := s.getJSON(ctx, "/api/v1/observability/overview", map[string]string{
		"from":          formatTime(from),
		"to":            formatTime(to),
		"tenantId":      tenantID,
		"periodSeconds": fmt.Sprintf("%d", periodSeconds),
	}, &result)
	return result, err
}

func (s *ObservabilityService) GetInventory(ctx context.Context, from, to *time.Time, tenantID string) ([]dto.ObservabilityInventoryItemResponse, error) {
	var result []dto.ObservabilityInventoryItemResponse
	err := s.getJSON(ctx, "/api/v1/observability/integrations/inventory", map[string]string{
		"from":     formatTime(from),
		"to":       formatTime(to),
		"tenantId": tenantID,
	}, &result)
	if result == nil {
		result = []dto.ObservabilityInventoryItemResponse{}
	}
	return result, err
}

func (s *ObservabilityService) GetAdminOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID, connectorID string) (dto.AdminObservabilityOverviewResponse, error) {
	var result dto.AdminObservabilityOverviewResponse
	err := s.getJSON(ctx, "/api/v1/admin/observability/overview", map[string]string{
		"from":          formatTime(from),
		"to":            formatTime(to),
		"tenantId":      tenantID,
		"connectorId":   connectorID,
		"periodSeconds": fmt.Sprintf("%d", periodSeconds),
	}, &result)
	return result, err
}

func (s *ObservabilityService) GetAdminTenantHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.TenantHealthListResponse, error) {
	var result dto.TenantHealthListResponse
	err := s.getJSON(ctx, "/api/v1/admin/observability/tenants", map[string]string{
		"from":   formatTime(from),
		"to":     formatTime(to),
		"status": status,
		"page":   fmt.Sprintf("%d", page),
		"size":   fmt.Sprintf("%d", size),
	}, &result)
	return result, err
}

func (s *ObservabilityService) GetAdminConnectorHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.ConnectorHealthListResponse, error) {
	var result dto.ConnectorHealthListResponse
	err := s.getJSON(ctx, "/api/v1/admin/observability/connectors", map[string]string{
		"from":   formatTime(from),
		"to":     formatTime(to),
		"status": status,
		"page":   fmt.Sprintf("%d", page),
		"size":   fmt.Sprintf("%d", size),
	}, &result)
	return result, err
}

func (s *ObservabilityService) getJSON(ctx context.Context, path string, query map[string]string, target any) error {
	endpoint, err := url.Parse(s.baseURL + path)
	if err != nil {
		return err
	}
	values := endpoint.Query()
	for key, value := range query {
		if value != "" {
			values.Set(key, value)
		}
	}
	endpoint.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("observability upstream returned %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func formatTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
