package httpx

import (
	"encoding/json"
	"net/http"
	"os"

	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/controller"
	"plexus-bff-service-go/internal/app/security"
	"plexus-bff-service-go/internal/app/tracing"
)

func NewRouter(
	cfg *config.Config,
	authController *controller.AuthController,
	connectorsController *controller.ConnectorsController,
	featureController *controller.FeatureController,
	featureFlagsController *controller.FeatureFlagsController,
	observabilityController *controller.ObservabilityController,
	adminObservabilityController *controller.AdminObservabilityController,
	testSessionController *controller.TestSessionController,
	adminUsersController *controller.AdminUsersController,
	adminGroupsController *controller.AdminGroupsController,
	adminPermissionsController *controller.AdminPermissionsController,
	domainsController *controller.DomainsController,
) (http.Handler, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{"status": "UP"}
		if cfg.Management.Health.Redis.Enabled {
			payload["redis"] = map[string]string{"status": "UNKNOWN"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	})
	mux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/swagger-ui.html")
	})
	mux.HandleFunc("/v3/api-docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, "static/openapi.yaml")
	})
	mux.HandleFunc("/auth/login", authController.Login)
	mux.HandleFunc("/auth/callback", authController.Callback)

	mux.Handle("/api/v1/features", protected(cfg, http.HandlerFunc(featureFlagsController.RetrieveFeatures)))
	mux.Handle("/api/v2/status", protected(cfg, http.HandlerFunc(featureController.Status)))
	mux.Handle("/api/v3/connectors", protected(cfg, http.HandlerFunc(connectorsController.Status)))
	mux.Handle("/api/v3/observability/overview", protected(cfg, http.HandlerFunc(observabilityController.Overview)))
	mux.Handle("/api/v3/observability/integrations/inventory", protected(cfg, http.HandlerFunc(observabilityController.Inventory)))
	mux.Handle("/api/v3/admin/observability/overview", protected(cfg, http.HandlerFunc(adminObservabilityController.Overview)))
	mux.Handle("/api/v3/admin/observability/tenants", protected(cfg, http.HandlerFunc(adminObservabilityController.Tenants)))
	mux.Handle("/api/v3/admin/observability/connectors", protected(cfg, http.HandlerFunc(adminObservabilityController.Connectors)))
	mux.Handle("GET /api/admin/users", protected(cfg, http.HandlerFunc(adminUsersController.List)))
	mux.Handle("POST /api/admin/users", protected(cfg, http.HandlerFunc(adminUsersController.Create)))
	mux.Handle("GET /api/admin/users/{id}", protected(cfg, http.HandlerFunc(adminUsersController.Get)))
	mux.Handle("PUT /api/admin/users/{id}", protected(cfg, http.HandlerFunc(adminUsersController.Update)))
	mux.Handle("DELETE /api/admin/users/{id}", protected(cfg, http.HandlerFunc(adminUsersController.Delete)))
	mux.Handle("GET /api/admin/groups", protected(cfg, http.HandlerFunc(adminGroupsController.List)))
	mux.Handle("POST /api/admin/groups", protected(cfg, http.HandlerFunc(adminGroupsController.Create)))
	mux.Handle("GET /api/admin/groups/{id}", protected(cfg, http.HandlerFunc(adminGroupsController.Get)))
	mux.Handle("PUT /api/admin/groups/{id}", protected(cfg, http.HandlerFunc(adminGroupsController.Update)))
	mux.Handle("DELETE /api/admin/groups/{id}", protected(cfg, http.HandlerFunc(adminGroupsController.Delete)))
	mux.Handle("GET /api/admin/permissions", protected(cfg, http.HandlerFunc(adminPermissionsController.List)))
	mux.Handle("POST /api/admin/permissions", protected(cfg, http.HandlerFunc(adminPermissionsController.Create)))
	mux.Handle("GET /api/admin/permissions/{id}", protected(cfg, http.HandlerFunc(adminPermissionsController.Get)))
	mux.Handle("PUT /api/admin/permissions/{id}", protected(cfg, http.HandlerFunc(adminPermissionsController.Update)))
	mux.Handle("DELETE /api/admin/permissions/{id}", protected(cfg, http.HandlerFunc(adminPermissionsController.Delete)))
	mux.Handle("GET /api/admin/domains/workspace", protected(cfg, http.HandlerFunc(domainsController.Workspace)))
	mux.Handle("GET /api/admin/domains", protected(cfg, http.HandlerFunc(domainsController.List)))
	mux.Handle("POST /api/admin/domains", protected(cfg, http.HandlerFunc(domainsController.Create)))
	mux.Handle("GET /api/admin/domains/approved", protected(cfg, http.HandlerFunc(domainsController.Approved)))
	mux.Handle("GET /api/admin/domains/{id}", protected(cfg, http.HandlerFunc(domainsController.Get)))
	mux.Handle("PUT /api/admin/domains/{id}", protected(cfg, http.HandlerFunc(domainsController.Update)))
	mux.Handle("DELETE /api/admin/domains/{id}", protected(cfg, http.HandlerFunc(domainsController.Delete)))
	mux.Handle("POST /api/admin/domains/{id}/decision", protected(cfg, http.HandlerFunc(domainsController.Decision)))

	if cfg.ActiveProfile() == "dev" || cfg.ActiveProfile() == "local-redis" {
		mux.Handle("/_test/session", protected(cfg, sessionWriter(http.HandlerFunc(testSessionController.CreateSession))))
	}

	return tracing.Middleware(mux), nil
}

func protected(cfg *config.Config, next http.Handler) http.Handler {
	if !cfg.Security.Enabled {
		return next
	}
	validator := security.NewValidator(cfg.Security.JWT)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if len(header) < len("Bearer ")+1 || header[:7] != "Bearer " {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		auth, err := validator.Parse(header[7:])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(security.WithAuthentication(r.Context(), auth)))
	})
}

func sessionWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = os.MkdirAll("target", 0o755)
		_ = os.WriteFile("target/session-marker.txt", []byte("test-key=test-value"), 0o644)
		next.ServeHTTP(w, r)
	})
}
