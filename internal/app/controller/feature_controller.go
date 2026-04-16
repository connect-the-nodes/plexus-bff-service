package controller

import (
	"encoding/json"
	"net/http"

	"plexus-bff-service-go/internal/app/apiresponse"
	"plexus-bff-service-go/internal/app/service"
)

type FeatureController struct {
	service service.FeatureService
}

func NewFeatureController(featureService service.FeatureService) *FeatureController {
	return &FeatureController{service: featureService}
}

func (c *FeatureController) Status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, apiresponse.OK(c.service.GetStatus()))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
