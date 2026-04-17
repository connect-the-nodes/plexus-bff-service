package impl

import (
	"context"
	"strings"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/common"
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
	"plexus-bff-service-go/internal/app/service"
)

type GroupAdminService struct {
	groups      repository.GroupRepository
	permissions repository.PermissionRepository
	mapper      *mapper.AdminMapper
}

func NewGroupAdminService(groups repository.GroupRepository, permissions repository.PermissionRepository) service.GroupAdminService {
	return &GroupAdminService{groups: groups, permissions: permissions, mapper: mapper.NewAdminMapper()}
}

func (s *GroupAdminService) ListGroups(ctx context.Context) ([]dto.PortalGroup, error) {
	items, err := s.groups.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToPortalGroups(items), nil
}

func (s *GroupAdminService) GetGroup(ctx context.Context, id string) (dto.PortalGroup, error) {
	item, err := s.groups.GetGroup(ctx, id)
	if err != nil {
		return dto.PortalGroup{}, err
	}
	return s.mapper.ToPortalGroup(item), nil
}

func (s *GroupAdminService) CreateGroup(ctx context.Context, input dto.CreatePortalGroupInput) (dto.PortalGroup, error) {
	group, err := s.buildGroup(ctx, common.NewID(), input)
	if err != nil {
		return dto.PortalGroup{}, err
	}
	item, err := s.groups.CreateGroup(ctx, group)
	if err != nil {
		return dto.PortalGroup{}, err
	}
	return s.mapper.ToPortalGroup(item), nil
}

func (s *GroupAdminService) UpdateGroup(ctx context.Context, id string, input dto.CreatePortalGroupInput) (dto.PortalGroup, error) {
	group, err := s.buildGroup(ctx, id, input)
	if err != nil {
		return dto.PortalGroup{}, err
	}
	item, err := s.groups.UpdateGroup(ctx, group)
	if err != nil {
		return dto.PortalGroup{}, err
	}
	return s.mapper.ToPortalGroup(item), nil
}

func (s *GroupAdminService) DeleteGroup(ctx context.Context, id string) error {
	return s.groups.DeleteGroup(ctx, id)
}

func (s *GroupAdminService) buildGroup(ctx context.Context, id string, input dto.CreatePortalGroupInput) (model.PortalGroup, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return model.PortalGroup{}, apperrors.NewValidation("name is required")
	}
	for _, permissionID := range input.PermissionIDs {
		if _, err := s.permissions.GetPermission(ctx, permissionID); err != nil {
			return model.PortalGroup{}, apperrors.NewValidation("one or more permissionIds do not exist")
		}
	}
	return model.PortalGroup{
		ID:            id,
		Name:          name,
		Description:   strings.TrimSpace(input.Description),
		PermissionIDs: append([]string(nil), input.PermissionIDs...),
	}, nil
}
