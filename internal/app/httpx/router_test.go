package httpx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/controller"
	"plexus-bff-service-go/internal/app/dto"
	repositoryimpl "plexus-bff-service-go/internal/app/repository/impl"
	serviceimpl "plexus-bff-service-go/internal/app/service/impl"
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

func TestAdminResourcesCRUDFlow(t *testing.T) {
	router := buildRouter(t, false)

	permissionID := createResource(t, router, http.MethodPost, "/api/admin/permissions", `{"name":"MANAGE_DOMAINS","description":"Manage domains"}`, http.StatusCreated)["id"].(string)
	groupPayload := createResource(t, router, http.MethodPost, "/api/admin/groups", `{"name":"Platform Admins","permissionIds":["`+permissionID+`"]}`, http.StatusCreated)
	groupID := groupPayload["id"].(string)
	userPayload := createResource(t, router, http.MethodPost, "/api/admin/users", `{"username":"suresh","password":"secret","email":"suresh@example.com","groupIds":["`+groupID+`"],"active":true}`, http.StatusCreated)
	if userPayload["username"] != "suresh" {
		t.Fatalf("expected created user to be returned")
	}

	domainPayload := createResource(t, router, http.MethodPost, "/api/admin/domains", `{"mode":"DRAFT","domain":{"name":"payments","ownerGroupId":"`+groupID+`","metadata":{"env":"dev"}}}`, http.StatusCreated)
	if domainPayload["status"] != "DRAFT" {
		t.Fatalf("expected domain status DRAFT, got %v", domainPayload["status"])
	}

	request := httptest.NewRequest(http.MethodGet, "/api/admin/domains/workspace", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected workspace status 200, got %d", recorder.Code)
	}

	var workspace map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &workspace); err != nil {
		t.Fatalf("decode workspace: %v", err)
	}
	if len(workspace["ownerGroups"].([]any)) != 1 {
		t.Fatalf("expected one owner group")
	}
	if len(workspace["domains"].([]any)) != 1 {
		t.Fatalf("expected one domain")
	}
}

func TestDeletePendingReviewDomainRejected(t *testing.T) {
	router := buildRouter(t, false)
	groupID := createResource(t, router, http.MethodPost, "/api/admin/groups", `{"name":"Approvers"}`, http.StatusCreated)["id"].(string)
	domainID := createResource(t, router, http.MethodPost, "/api/admin/domains", `{"mode":"SUBMIT","domain":{"name":"identity","ownerGroupId":"`+groupID+`"}}`, http.StatusCreated)["id"].(string)

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/domains/"+domainID, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected delete to fail with 400, got %d", recorder.Code)
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

	adminRepo := repositoryimpl.NewMemoryAdminRepository()

	router, err := NewRouter(
		cfg,
		controller.NewAuthController(config.CognitoConfig{}),
		controller.NewConnectorsController(&testConnectorsService{}),
		controller.NewFeatureController(&testFeatureService{}),
		controller.NewFeatureFlagsController(nil),
		controller.NewObservabilityController(&testObservabilityService{}),
		controller.NewAdminObservabilityController(&testObservabilityService{}),
		controller.NewTestSessionController(),
		controller.NewAdminUsersController(serviceimpl.NewUserAdminService(adminRepo.Users(), adminRepo.Groups())),
		controller.NewAdminGroupsController(serviceimpl.NewGroupAdminService(adminRepo.Groups(), adminRepo.Permissions())),
		controller.NewAdminPermissionsController(serviceimpl.NewPermissionAdminService(adminRepo.Permissions())),
		controller.NewDomainsController(serviceimpl.NewDomainAdminService(adminRepo.Domains(), adminRepo.Groups())),
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

func createResource(t *testing.T, router http.Handler, method, path, body string, expectedStatus int) map[string]any {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != expectedStatus {
		t.Fatalf("expected %d for %s %s, got %d and body %s", expectedStatus, method, path, recorder.Code, recorder.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response for %s %s: %v", method, path, err)
	}
	return payload
}
