package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/apiresponse"
	"plexus-bff-service-go/internal/app/security"
	"plexus-bff-service-go/internal/app/service"
)

type AdminObservabilityController struct {
	service service.ObservabilityService
}

func NewAdminObservabilityController(observabilityService service.ObservabilityService) *AdminObservabilityController {
	return &AdminObservabilityController{service: observabilityService}
}

func (c *AdminObservabilityController) Overview(w http.ResponseWriter, r *http.Request) {
	if !ensureAdmin(r) {
		writeJSON(w, http.StatusForbidden, apiresponse.Error("Admin role required"))
		return
	}
	from, to := parseTimeRange(r)
	response, err := c.service.GetAdminOverview(
		r.Context(),
		from,
		to,
		parseIntDefault(r.URL.Query().Get("periodSeconds"), 300),
		r.URL.Query().Get("tenantId"),
		r.URL.Query().Get("connectorId"),
	)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}

func (c *AdminObservabilityController) Tenants(w http.ResponseWriter, r *http.Request) {
	if !ensureAdmin(r) {
		writeJSON(w, http.StatusForbidden, apiresponse.Error("Admin role required"))
		return
	}
	from, to := parseTimeRange(r)
	response, err := c.service.GetAdminTenantHealth(
		r.Context(),
		from,
		to,
		r.URL.Query().Get("status"),
		parseIntDefault(r.URL.Query().Get("page"), 1),
		parseIntDefault(r.URL.Query().Get("size"), 50),
	)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}

func (c *AdminObservabilityController) Connectors(w http.ResponseWriter, r *http.Request) {
	if !ensureAdmin(r) {
		writeJSON(w, http.StatusForbidden, apiresponse.Error("Admin role required"))
		return
	}
	from, to := parseTimeRange(r)
	response, err := c.service.GetAdminConnectorHealth(
		r.Context(),
		from,
		to,
		r.URL.Query().Get("status"),
		parseIntDefault(r.URL.Query().Get("page"), 1),
		parseIntDefault(r.URL.Query().Get("size"), 50),
	)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}

func ensureAdmin(r *http.Request) bool {
	auth := security.FromContext(r.Context())
	if auth == nil {
		return true
	}
	adminAuthorities := map[string]struct{}{
		"ROLE_ADMIN":      {},
		"ROLE_AXIS_ADMIN": {},
		"ROLE_SUPER_ADMIN": {},
		"SCOPE_admin":     {},
		"SCOPE_axis.admin": {},
	}
	for _, authority := range auth.Authorities {
		if _, ok := adminAuthorities[authority]; ok {
			return true
		}
	}
	return false
}
