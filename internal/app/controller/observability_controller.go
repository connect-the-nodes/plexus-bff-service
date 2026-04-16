package controller

import (
	"net/http"
	"time"

	"plexus-bff-service-go/internal/app/apiresponse"
	"plexus-bff-service-go/internal/app/security"
	"plexus-bff-service-go/internal/app/service"
)

type ObservabilityController struct {
	service service.ObservabilityService
}

func NewObservabilityController(observabilityService service.ObservabilityService) *ObservabilityController {
	return &ObservabilityController{service: observabilityService}
}

func (c *ObservabilityController) Overview(w http.ResponseWriter, r *http.Request) {
	from, to := parseTimeRange(r)
	tenantID, status, message := resolveTenantID(r)
	if status != http.StatusOK {
		writeJSON(w, status, apiresponse.Error(message))
		return
	}

	periodSeconds := parseIntDefault(r.URL.Query().Get("periodSeconds"), 300)
	response, err := c.service.GetOverview(r.Context(), from, to, periodSeconds, tenantID)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}

func (c *ObservabilityController) Inventory(w http.ResponseWriter, r *http.Request) {
	from, to := parseTimeRange(r)
	tenantID, status, message := resolveTenantID(r)
	if status != http.StatusOK {
		writeJSON(w, status, apiresponse.Error(message))
		return
	}

	response, err := c.service.GetInventory(r.Context(), from, to, tenantID)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}

func resolveTenantID(r *http.Request) (string, int, string) {
	auth := security.FromContext(r.Context())
	queryTenant := r.URL.Query().Get("tenantId")

	if auth != nil {
		jwtTenant := getTenantFromClaims(auth.Claims)
		if jwtTenant == "" {
			return "", http.StatusForbidden, "Tenant claim missing in JWT"
		}
		if queryTenant != "" && queryTenant != jwtTenant {
			return "", http.StatusForbidden, "Tenant mismatch"
		}
		return jwtTenant, http.StatusOK, ""
	}

	if queryTenant != "" {
		return queryTenant, http.StatusOK, ""
	}
	return "default", http.StatusOK, ""
}

func getTenantFromClaims(claims map[string]any) string {
	for _, key := range []string{"tenantId", "custom:tenantId", "tenant_id"} {
		if value, ok := claims[key].(string); ok && value != "" {
			return value
		}
	}
	if groups, ok := claims["cognito:groups"].([]any); ok {
		for _, group := range groups {
			text, ok := group.(string)
			if !ok || text == "" {
				continue
			}
			switch {
			case len(text) > len("TENANT_") && text[:len("TENANT_")] == "TENANT_":
				return text[len("TENANT_"):]
			case len(text) > len("tenant:") && text[:len("tenant:")] == "tenant:":
				return text[len("tenant:"):]
			case len(text) > len("TENANT#") && text[:len("TENANT#")] == "TENANT#":
				return text[len("TENANT#"):]
			}
		}
	}
	return ""
}

func parseTimeRange(r *http.Request) (*time.Time, *time.Time) {
	var from *time.Time
	var to *time.Time
	if raw := r.URL.Query().Get("from"); raw != "" {
		if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
			from = &parsed
		}
	}
	if raw := r.URL.Query().Get("to"); raw != "" {
		if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
			to = &parsed
		}
	}
	return from, to
}
