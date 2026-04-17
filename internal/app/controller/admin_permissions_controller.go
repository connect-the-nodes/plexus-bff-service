package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type AdminPermissionsController struct {
	service service.PermissionAdminService
}

func NewAdminPermissionsController(permissionService service.PermissionAdminService) *AdminPermissionsController {
	return &AdminPermissionsController{service: permissionService}
}

func (c *AdminPermissionsController) List(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListPermissions(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, items)
}

func (c *AdminPermissionsController) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePermissionInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.CreatePermission(r.Context(), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusCreated, item)
}

func (c *AdminPermissionsController) Get(w http.ResponseWriter, r *http.Request) {
	item, err := c.service.GetPermission(r.Context(), r.PathValue("id"))
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminPermissionsController) Update(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePermissionInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.UpdatePermission(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminPermissionsController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeletePermission(r.Context(), r.PathValue("id")); err != nil {
		writeOpenAPIError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
