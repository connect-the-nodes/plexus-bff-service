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

type PermissionAdminService struct {
	repository repository.PermissionRepository
	mapper     *mapper.AdminMapper
}

func NewPermissionAdminService(permissionRepository repository.PermissionRepository) service.PermissionAdminService {
	return &PermissionAdminService{repository: permissionRepository, mapper: mapper.NewAdminMapper()}
}

func (s *PermissionAdminService) ListPermissions(ctx context.Context) ([]dto.PermissionDefinition, error) {
	items, err := s.repository.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToPermissions(items), nil
}

func (s *PermissionAdminService) GetPermission(ctx context.Context, id string) (dto.PermissionDefinition, error) {
	item, err := s.repository.GetPermission(ctx, id)
	if err != nil {
		return dto.PermissionDefinition{}, err
	}
	return s.mapper.ToPermission(item), nil
}

func (s *PermissionAdminService) CreatePermission(ctx context.Context, input dto.CreatePermissionInput) (dto.PermissionDefinition, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return dto.PermissionDefinition{}, apperrors.NewValidation("name is required")
	}
	item, err := s.repository.CreatePermission(ctx, model.PermissionDefinition{
		ID:          common.NewID(),
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	})
	if err != nil {
		return dto.PermissionDefinition{}, err
	}
	return s.mapper.ToPermission(item), nil
}

func (s *PermissionAdminService) UpdatePermission(ctx context.Context, id string, input dto.CreatePermissionInput) (dto.PermissionDefinition, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return dto.PermissionDefinition{}, apperrors.NewValidation("name is required")
	}
	item, err := s.repository.UpdatePermission(ctx, model.PermissionDefinition{
		ID:          id,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	})
	if err != nil {
		return dto.PermissionDefinition{}, err
	}
	return s.mapper.ToPermission(item), nil
}

func (s *PermissionAdminService) DeletePermission(ctx context.Context, id string) error {
	return s.repository.DeletePermission(ctx, id)
}
