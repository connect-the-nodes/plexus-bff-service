package httpx

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/controller"
	"plexus-bff-service-go/internal/app/dto"
)

type testFeatureService struct{}

func (s *testFeatureService) GetStatus() dto.FeatureResponse {
	return dto.FeatureResponse{Status: "OK", Message: "Service is running"}
}

type testConnectorsService struct{}

func (s *testConnectorsService) GetConnectors(ctx context.Context) (dto.ConnectorsResponse, error) {
	return dto.ConnectorsResponse{Status: "OK", Message: "Connectors Service is running"}, nil
}

type testObservabilityService struct{}

func (s *testObservabilityService) GetOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID string) (dto.ObservabilityOverviewResponse, error) {
	return dto.ObservabilityOverviewResponse{}, nil
}

func (s *testObservabilityService) GetInventory(ctx context.Context, from, to *time.Time, tenantID string) ([]dto.ObservabilityInventoryItemResponse, error) {
	return []dto.ObservabilityInventoryItemResponse{}, nil
}

func (s *testObservabilityService) GetAdminOverview(ctx context.Context, from, to *time.Time, periodSeconds int, tenantID, connectorID string) (dto.AdminObservabilityOverviewResponse, error) {
	return dto.AdminObservabilityOverviewResponse{}, nil
}

func (s *testObservabilityService) GetAdminTenantHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.TenantHealthListResponse, error) {
	return dto.TenantHealthListResponse{}, nil
}

func (s *testObservabilityService) GetAdminConnectorHealth(ctx context.Context, from, to *time.Time, status string, page, size int) (dto.ConnectorHealthListResponse, error) {
	return dto.ConnectorHealthListResponse{}, nil
}

func TestHealthEndpointIsPermittedWithoutAuth(t *testing.T) {
	router := buildRouter(t, true)

	request := httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestProtectedEndpointBlockedWithoutAuth(t *testing.T) {
	router := buildRouter(t, true)

	request := httptest.NewRequest(http.MethodGet, "/api/v2/status", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestProtectedEndpointAllowedWithJWT(t *testing.T) {
	router := buildRouter(t, true)

	request := httptest.NewRequest(http.MethodGet, "/api/v2/status", nil)
	request.Header.Set("Authorization", "Bearer "+fakeJWT(`{"cognito:groups":["ADMIN"]}`))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestCorrelationIDEchoedOrGenerated(t *testing.T) {
	router := buildRouter(t, false)

	request := httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	request.Header.Set("X-Correlation-Id", "test-correlation-id")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Header().Get("X-Correlation-Id") != "test-correlation-id" {
		t.Fatalf("expected correlation id to be echoed")
	}

	request = httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Header().Get("X-Correlation-Id") == "" {
		t.Fatalf("expected correlation id to be generated")
	}
}

func buildRouter(t *testing.T, securityEnabled bool) http.Handler {
	t.Helper()
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Enabled: securityEnabled,
			JWT: config.SecurityJWTConfig{
				AuthoritiesClaim: "cognito:groups",
				RolePrefix:       "ROLE_",
			},
		},
		Observability: config.ObservabilityConfig{
			Service: config.ObservabilityServiceConfig{BaseURL: "http://localhost:8081"},
		},
	}

	router, err := NewRouter(
		cfg,
		controller.NewAuthController(config.CognitoConfig{}),
		controller.NewConnectorsController(&testConnectorsService{}),
		controller.NewFeatureController(&testFeatureService{}),
		controller.NewFeatureFlagsController(nil),
		controller.NewObservabilityController(&testObservabilityService{}),
		controller.NewAdminObservabilityController(&testObservabilityService{}),
		controller.NewTestSessionController(),
	)
	if err != nil {
		t.Fatalf("create router: %v", err)
	}
	return router
}

func fakeJWT(payload string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	body := base64.RawURLEncoding.EncodeToString([]byte(payload))
	return strings.Join([]string{header, body, ""}, ".")
}
