package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type AdminGroupsController struct {
	service service.GroupAdminService
}

func NewAdminGroupsController(groupService service.GroupAdminService) *AdminGroupsController {
	return &AdminGroupsController{service: groupService}
}

func (c *AdminGroupsController) List(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListGroups(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, items)
}

func (c *AdminGroupsController) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePortalGroupInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.CreateGroup(r.Context(), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusCreated, item)
}

func (c *AdminGroupsController) Get(w http.ResponseWriter, r *http.Request) {
	item, err := c.service.GetGroup(r.Context(), r.PathValue("id"))
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminGroupsController) Update(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePortalGroupInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.UpdateGroup(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminGroupsController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteGroup(r.Context(), r.PathValue("id")); err != nil {
		writeOpenAPIError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
