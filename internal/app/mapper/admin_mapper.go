package mapper

import (
	"time"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/model"
)

type AdminMapper struct{}

func NewAdminMapper() *AdminMapper {
	return &AdminMapper{}
}

func (m *AdminMapper) ToPortalUser(modelUser model.PortalUser) dto.PortalUser {
	return dto.PortalUser{
		ID:          modelUser.ID,
		Username:    modelUser.Username,
		Email:       modelUser.Email,
		DisplayName: modelUser.DisplayName,
		Role:        modelUser.Role,
		Active:      modelUser.Active,
		GroupIDs:    append([]string(nil), modelUser.GroupIDs...),
	}
}

func (m *AdminMapper) ToPortalUsers(users []model.PortalUser) []dto.PortalUser {
	result := make([]dto.PortalUser, 0, len(users))
	for _, user := range users {
		result = append(result, m.ToPortalUser(user))
	}
	return result
}

func (m *AdminMapper) ToPermission(modelPermission model.PermissionDefinition) dto.PermissionDefinition {
	return dto.PermissionDefinition{
		ID:          modelPermission.ID,
		Name:        modelPermission.Name,
		Description: modelPermission.Description,
	}
}

func (m *AdminMapper) ToPermissions(items []model.PermissionDefinition) []dto.PermissionDefinition {
	result := make([]dto.PermissionDefinition, 0, len(items))
	for _, item := range items {
		result = append(result, m.ToPermission(item))
	}
	return result
}

func (m *AdminMapper) ToPortalGroup(modelGroup model.PortalGroup) dto.PortalGroup {
	return dto.PortalGroup{
		ID:            modelGroup.ID,
		Name:          modelGroup.Name,
		Description:   modelGroup.Description,
		PermissionIDs: append([]string(nil), modelGroup.PermissionIDs...),
	}
}

func (m *AdminMapper) ToPortalGroups(groups []model.PortalGroup) []dto.PortalGroup {
	result := make([]dto.PortalGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, m.ToPortalGroup(group))
	}
	return result
}

func (m *AdminMapper) ToRegisteredDomain(domain model.RegisteredDomain) dto.RegisteredDomain {
	comments := make([]dto.ReviewComment, 0, len(domain.Review))
	for _, review := range domain.Review {
		comments = append(comments, dto.ReviewComment{
			Author:    review.Author,
			AuthorID:  review.AuthorID,
			Decision:  review.Decision,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	return dto.RegisteredDomain{
		ID:           domain.ID,
		Name:         domain.Name,
		Description:  domain.Description,
		OwnerGroupID: domain.OwnerGroupID,
		Status:       domain.Status,
		Review:       comments,
		Metadata:     cloneMap(domain.Metadata),
	}
}

func (m *AdminMapper) ToRegisteredDomains(items []model.RegisteredDomain) []dto.RegisteredDomain {
	result := make([]dto.RegisteredDomain, 0, len(items))
	for _, item := range items {
		result = append(result, m.ToRegisteredDomain(item))
	}
	return result
}

func cloneMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
