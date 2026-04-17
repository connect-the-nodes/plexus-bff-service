package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type AdminUsersController struct {
	service service.UserAdminService
}

func NewAdminUsersController(userService service.UserAdminService) *AdminUsersController {
	return &AdminUsersController{service: userService}
}

func (c *AdminUsersController) List(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListUsers(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, items)
}

func (c *AdminUsersController) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePortalUserInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.CreateUser(r.Context(), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusCreated, item)
}

func (c *AdminUsersController) Get(w http.ResponseWriter, r *http.Request) {
	item, err := c.service.GetUser(r.Context(), r.PathValue("id"))
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminUsersController) Update(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePortalUserInput
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.UpdateUser(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *AdminUsersController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteUser(r.Context(), r.PathValue("id")); err != nil {
		writeOpenAPIError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
