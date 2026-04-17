package impl

import (
	"context"
	"sort"
	"sync"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
)

type MemoryAdminRepository struct {
	mu          sync.RWMutex
	users       map[string]model.PortalUser
	permissions map[string]model.PermissionDefinition
	groups      map[string]model.PortalGroup
	domains     map[string]model.RegisteredDomain
}

func NewMemoryAdminRepository() *MemoryAdminRepository {
	return &MemoryAdminRepository{
		users:       map[string]model.PortalUser{},
		permissions: map[string]model.PermissionDefinition{},
		groups:      map[string]model.PortalGroup{},
		domains:     map[string]model.RegisteredDomain{},
	}
}

func (r *MemoryAdminRepository) Users() repository.UserRepository             { return r }
func (r *MemoryAdminRepository) Permissions() repository.PermissionRepository { return r }
func (r *MemoryAdminRepository) Groups() repository.GroupRepository           { return r }
func (r *MemoryAdminRepository) Domains() repository.DomainRepository         { return r }

func (r *MemoryAdminRepository) ListUsers(context.Context) ([]model.PortalUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.PortalUser, 0, len(r.users))
	for _, item := range r.users {
		items = append(items, cloneUser(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Username < items[j].Username })
	return items, nil
}

func (r *MemoryAdminRepository) GetUser(_ context.Context, id string) (model.PortalUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.users[id]
	if !ok {
		return model.PortalUser{}, apperrors.NewNotFound("user not found")
	}
	return cloneUser(item), nil
}

func (r *MemoryAdminRepository) CreateUser(_ context.Context, user model.PortalUser) (model.PortalUser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.users {
		if existing.Username == user.Username {
			return model.PortalUser{}, apperrors.NewConflict("username already exists")
		}
		if user.Email != "" && existing.Email == user.Email {
			return model.PortalUser{}, apperrors.NewConflict("email already exists")
		}
	}
	r.users[user.ID] = cloneUser(user)
	return cloneUser(user), nil
}

func (r *MemoryAdminRepository) UpdateUser(_ context.Context, user model.PortalUser) (model.PortalUser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return model.PortalUser{}, apperrors.NewNotFound("user not found")
	}
	for id, existing := range r.users {
		if id == user.ID {
			continue
		}
		if existing.Username == user.Username {
			return model.PortalUser{}, apperrors.NewConflict("username already exists")
		}
		if user.Email != "" && existing.Email == user.Email {
			return model.PortalUser{}, apperrors.NewConflict("email already exists")
		}
	}
	r.users[user.ID] = cloneUser(user)
	return cloneUser(user), nil
}

func (r *MemoryAdminRepository) DeleteUser(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return apperrors.NewNotFound("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *MemoryAdminRepository) ListPermissions(context.Context) ([]model.PermissionDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.PermissionDefinition, 0, len(r.permissions))
	for _, item := range r.permissions {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *MemoryAdminRepository) GetPermission(_ context.Context, id string) (model.PermissionDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.permissions[id]
	if !ok {
		return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
	}
	return item, nil
}

func (r *MemoryAdminRepository) CreatePermission(_ context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.permissions {
		if existing.Name == permission.Name {
			return model.PermissionDefinition{}, apperrors.NewConflict("permission already exists")
		}
	}
	r.permissions[permission.ID] = permission
	return permission, nil
}

func (r *MemoryAdminRepository) UpdatePermission(_ context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[permission.ID]; !ok {
		return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
	}
	for id, existing := range r.permissions {
		if id != permission.ID && existing.Name == permission.Name {
			return model.PermissionDefinition{}, apperrors.NewConflict("permission already exists")
		}
	}
	r.permissions[permission.ID] = permission
	return permission, nil
}

func (r *MemoryAdminRepository) DeletePermission(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[id]; !ok {
		return apperrors.NewNotFound("permission not found")
	}
	delete(r.permissions, id)
	for groupID, group := range r.groups {
		group.PermissionIDs = removeID(group.PermissionIDs, id)
		r.groups[groupID] = group
	}
	return nil
}

func (r *MemoryAdminRepository) ListGroups(context.Context) ([]model.PortalGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.PortalGroup, 0, len(r.groups))
	for _, item := range r.groups {
		items = append(items, cloneGroup(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *MemoryAdminRepository) GetGroup(_ context.Context, id string) (model.PortalGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.groups[id]
	if !ok {
		return model.PortalGroup{}, apperrors.NewNotFound("group not found")
	}
	return cloneGroup(item), nil
}

func (r *MemoryAdminRepository) CreateGroup(_ context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.groups {
		if existing.Name == group.Name {
			return model.PortalGroup{}, apperrors.NewConflict("group already exists")
		}
	}
	r.groups[group.ID] = cloneGroup(group)
	return cloneGroup(group), nil
}

func (r *MemoryAdminRepository) UpdateGroup(_ context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.groups[group.ID]; !ok {
		return model.PortalGroup{}, apperrors.NewNotFound("group not found")
	}
	for id, existing := range r.groups {
		if id != group.ID && existing.Name == group.Name {
			return model.PortalGroup{}, apperrors.NewConflict("group already exists")
		}
	}
	r.groups[group.ID] = cloneGroup(group)
	return cloneGroup(group), nil
}

func (r *MemoryAdminRepository) DeleteGroup(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.groups[id]; !ok {
		return apperrors.NewNotFound("group not found")
	}
	for _, domain := range r.domains {
		if domain.OwnerGroupID == id {
			return apperrors.NewValidation("group is referenced by one or more domains")
		}
	}
	delete(r.groups, id)
	for userID, user := range r.users {
		user.GroupIDs = removeID(user.GroupIDs, id)
		r.users[userID] = user
	}
	return nil
}

func (r *MemoryAdminRepository) ListDomains(context.Context) ([]model.RegisteredDomain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.RegisteredDomain, 0, len(r.domains))
	for _, item := range r.domains {
		items = append(items, cloneDomain(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *MemoryAdminRepository) ListApprovedDomains(ctx context.Context) ([]model.RegisteredDomain, error) {
	items, err := r.ListDomains(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]model.RegisteredDomain, 0)
	for _, item := range items {
		if item.Status == "APPROVED" {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (r *MemoryAdminRepository) GetDomain(_ context.Context, id string) (model.RegisteredDomain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.domains[id]
	if !ok {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	return cloneDomain(item), nil
}

func (r *MemoryAdminRepository) CreateDomain(_ context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.domains {
		if existing.Name == domain.Name {
			return model.RegisteredDomain{}, apperrors.NewConflict("domain already exists")
		}
	}
	r.domains[domain.ID] = cloneDomain(domain)
	return cloneDomain(domain), nil
}

func (r *MemoryAdminRepository) UpdateDomain(_ context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.domains[domain.ID]; !ok {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	for id, existing := range r.domains {
		if id != domain.ID && existing.Name == domain.Name {
			return model.RegisteredDomain{}, apperrors.NewConflict("domain already exists")
		}
	}
	r.domains[domain.ID] = cloneDomain(domain)
	return cloneDomain(domain), nil
}

func (r *MemoryAdminRepository) DeleteDomain(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.domains[id]; !ok {
		return apperrors.NewNotFound("domain not found")
	}
	delete(r.domains, id)
	return nil
}

func (r *MemoryAdminRepository) AddDomainReview(_ context.Context, id string, comment model.ReviewComment, nextStatus string) (model.RegisteredDomain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.domains[id]
	if !ok {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	item.Review = append(item.Review, comment)
	item.Status = nextStatus
	r.domains[id] = item
	return cloneDomain(item), nil
}

func cloneUser(item model.PortalUser) model.PortalUser {
	item.GroupIDs = append([]string(nil), item.GroupIDs...)
	return item
}

func cloneGroup(item model.PortalGroup) model.PortalGroup {
	item.PermissionIDs = append([]string(nil), item.PermissionIDs...)
	return item
}

func cloneDomain(item model.RegisteredDomain) model.RegisteredDomain {
	item.Review = append([]model.ReviewComment(nil), item.Review...)
	item.Metadata = cloneMetadata(item.Metadata)
	return item
}

func cloneMetadata(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	result := make(map[string]any, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}

func removeID(items []string, target string) []string {
	filtered := make([]string, 0, len(items))
	for _, item := range items {
		if item != target {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
