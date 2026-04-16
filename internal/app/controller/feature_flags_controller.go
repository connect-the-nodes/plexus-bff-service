package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/apiresponse"
	"plexus-bff-service-go/internal/app/service"
)

type FeatureFlagsController struct {
	service *service.FeatureFlagsService
}

func NewFeatureFlagsController(featureFlagsService *service.FeatureFlagsService) *FeatureFlagsController {
	return &FeatureFlagsController{service: featureFlagsService}
}

func (c *FeatureFlagsController) RetrieveFeatures(w http.ResponseWriter, r *http.Request) {
	response, err := c.service.RetrieveFeatures(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}
