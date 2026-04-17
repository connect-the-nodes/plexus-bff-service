package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type PortalUser struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Password    string   `json:"password,omitempty"`
	Email       string   `json:"email,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Role        string   `json:"role"`
	Active      bool     `json:"active"`
	GroupIDs    []string `json:"groupIds"`
}

type CreatePortalUserInput struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Role        string   `json:"role"`
	Active      *bool    `json:"active"`
	GroupIDs    []string `json:"groupIds"`
}

type PermissionDefinition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreatePermissionInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PortalGroup struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	PermissionIDs []string `json:"permissionIds"`
}

type CreatePortalGroupInput struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds"`
}

type ReviewComment struct {
	Author    string `json:"author"`
	AuthorID  string `json:"authorId,omitempty"`
	Decision  string `json:"decision"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"`
}

type RegisteredDomain struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	OwnerGroupID string          `json:"ownerGroupId"`
	Status       string          `json:"status"`
	Review       []ReviewComment `json:"review"`
	Metadata     map[string]any  `json:"metadata,omitempty"`
}

type DomainWorkspace struct {
	OwnerGroups []PortalGroup      `json:"ownerGroups"`
	Domains     []RegisteredDomain `json:"domains"`
}

type DomainRegistrationInput struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	OwnerGroupID string         `json:"ownerGroupId"`
	Metadata     map[string]any `json:"metadata"`
}

type DomainMutationRequest struct {
	Mode   string                  `json:"mode"`
	Domain DomainRegistrationInput `json:"domain"`
}

type ReviewDecisionRequest struct {
	Author   string `json:"author"`
	AuthorID string `json:"authorId"`
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
}
