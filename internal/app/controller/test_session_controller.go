package controller

import (
	"net/http"
)

type TestSessionController struct{}

func NewTestSessionController() *TestSessionController {
	return &TestSessionController{}
}

func (c *TestSessionController) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
