package model

import "time"

type PortalUser struct {
	ID           string
	Username     string
	PasswordHash string
	Email        string
	DisplayName  string
	Role         string
	Active       bool
	GroupIDs     []string
}

type PermissionDefinition struct {
	ID          string
	Name        string
	Description string
}

type PortalGroup struct {
	ID            string
	Name          string
	Description   string
	PermissionIDs []string
}

type ReviewComment struct {
	ID        string
	Author    string
	AuthorID  string
	Decision  string
	Comment   string
	CreatedAt time.Time
}

type RegisteredDomain struct {
	ID           string
	Name         string
	Description  string
	OwnerGroupID string
	Status       string
	Review       []ReviewComment
	Metadata     map[string]any
}
