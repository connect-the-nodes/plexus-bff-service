package controller

import (
	"net/http"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type DomainsController struct {
	service service.DomainAdminService
}

func NewDomainsController(domainService service.DomainAdminService) *DomainsController {
	return &DomainsController{service: domainService}
}

func (c *DomainsController) Workspace(w http.ResponseWriter, r *http.Request) {
	item, err := c.service.GetWorkspace(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *DomainsController) List(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListDomains(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, items)
}

func (c *DomainsController) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.DomainMutationRequest
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.CreateDomain(r.Context(), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusCreated, item)
}

func (c *DomainsController) Approved(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListApprovedDomains(r.Context())
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, items)
}

func (c *DomainsController) Get(w http.ResponseWriter, r *http.Request) {
	item, err := c.service.GetDomain(r.Context(), r.PathValue("id"))
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *DomainsController) Update(w http.ResponseWriter, r *http.Request) {
	var input dto.DomainMutationRequest
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.UpdateDomain(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}

func (c *DomainsController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteDomain(r.Context(), r.PathValue("id")); err != nil {
		writeOpenAPIError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *DomainsController) Decision(w http.ResponseWriter, r *http.Request) {
	var input dto.ReviewDecisionRequest
	if !decodeOpenAPIJSON(w, r, &input) {
		return
	}
	item, err := c.service.ReviewDomain(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeOpenAPIError(w, err)
		return
	}
	writeOpenAPIJSON(w, http.StatusOK, item)
}
