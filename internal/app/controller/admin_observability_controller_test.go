package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"plexus-bff-service-go/internal/app/security"
)

func TestAdminObservabilityAllowsAdminRole(t *testing.T) {
	service := &fakeObservabilityService{}
	controller := NewAdminObservabilityController(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v3/admin/observability/overview", nil)
	request = request.WithContext(security.WithAuthentication(request.Context(), &security.Authentication{
		Authorities: []string{"ROLE_ADMIN"},
	}))
	recorder := httptest.NewRecorder()

	controller.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestAdminObservabilityRejectsNonAdminRole(t *testing.T) {
	service := &fakeObservabilityService{}
	controller := NewAdminObservabilityController(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v3/admin/observability/overview", nil)
	request = request.WithContext(security.WithAuthentication(request.Context(), &security.Authentication{
		Authorities: []string{"ROLE_TENANT_VIEWER"},
	}))
	recorder := httptest.NewRecorder()

	controller.Overview(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}
