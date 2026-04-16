package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"plexus-bff-service-go/internal/app/config"
)

type AuthController struct {
	properties config.CognitoConfig
}

func NewAuthController(properties config.CognitoConfig) *AuthController {
	return &AuthController{properties: properties}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if !c.properties.Enabled {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	scopes := c.properties.Scopes
	if scopes == "" {
		scopes = "openid,email,profile"
	}
	target := fmt.Sprintf("%s/oauth2/authorize?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s",
		c.properties.Domain,
		url.QueryEscape(c.properties.ClientID),
		url.QueryEscape(strings.Join(splitScopes(scopes), " ")),
		url.QueryEscape(c.properties.RedirectURI),
		url.QueryEscape("generated-state"),
	)
	http.Redirect(w, r, target, http.StatusFound)
}

func (c *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	if !c.properties.Enabled {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	target := fmt.Sprintf("%s?code=%s", c.properties.PostLoginRedirectURI, url.QueryEscape(code))
	if state != "" {
		target += "&state=" + url.QueryEscape(state)
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func splitScopes(scopes string) []string {
	parts := strings.Split(scopes, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return []string{"openid", "email", "profile"}
	}
	return result
}
