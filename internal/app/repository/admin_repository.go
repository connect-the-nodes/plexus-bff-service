package repository

import (
	"context"

	"plexus-bff-service-go/internal/app/model"
)

type UserRepository interface {
	ListUsers(ctx context.Context) ([]model.PortalUser, error)
	GetUser(ctx context.Context, id string) (model.PortalUser, error)
	CreateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error)
	UpdateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error)
	DeleteUser(ctx context.Context, id string) error
}

type PermissionRepository interface {
	ListPermissions(ctx context.Context) ([]model.PermissionDefinition, error)
	GetPermission(ctx context.Context, id string) (model.PermissionDefinition, error)
	CreatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error)
	UpdatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error)
	DeletePermission(ctx context.Context, id string) error
}

type GroupRepository interface {
	ListGroups(ctx context.Context) ([]model.PortalGroup, error)
	GetGroup(ctx context.Context, id string) (model.PortalGroup, error)
	CreateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error)
	UpdateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error)
	DeleteGroup(ctx context.Context, id string) error
}

type DomainRepository interface {
	ListDomains(ctx context.Context) ([]model.RegisteredDomain, error)
	ListApprovedDomains(ctx context.Context) ([]model.RegisteredDomain, error)
	GetDomain(ctx context.Context, id string) (model.RegisteredDomain, error)
	CreateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error)
	UpdateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error)
	DeleteDomain(ctx context.Context, id string) error
	AddDomainReview(ctx context.Context, id string, comment model.ReviewComment, nextStatus string) (model.RegisteredDomain, error)
}
