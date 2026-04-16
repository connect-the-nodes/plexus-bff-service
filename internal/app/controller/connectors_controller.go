package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/apiresponse"
	"plexus-bff-service-go/internal/app/service"
)

type ConnectorsController struct {
	service service.ConnectorsService
}

func NewConnectorsController(connectorsService service.ConnectorsService) *ConnectorsController {
	return &ConnectorsController{service: connectorsService}
}

func (c *ConnectorsController) Status(w http.ResponseWriter, r *http.Request) {
	response, err := c.service.GetConnectors(r.Context())
	if err != nil {
		writeJSON(w, http.StatusForbidden, apiresponse.Error(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, apiresponse.OK(response))
}
