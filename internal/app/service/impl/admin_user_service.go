package impl

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/common"
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
	"plexus-bff-service-go/internal/app/service"
)

type UserAdminService struct {
	users  repository.UserRepository
	groups repository.GroupRepository
	mapper *mapper.AdminMapper
}

func NewUserAdminService(users repository.UserRepository, groups repository.GroupRepository) service.UserAdminService {
	return &UserAdminService{users: users, groups: groups, mapper: mapper.NewAdminMapper()}
}

func (s *UserAdminService) ListUsers(ctx context.Context) ([]dto.PortalUser, error) {
	users, err := s.users.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToPortalUsers(users), nil
}

func (s *UserAdminService) GetUser(ctx context.Context, id string) (dto.PortalUser, error) {
	user, err := s.users.GetUser(ctx, id)
	if err != nil {
		return dto.PortalUser{}, err
	}
	return s.mapper.ToPortalUser(user), nil
}

func (s *UserAdminService) CreateUser(ctx context.Context, input dto.CreatePortalUserInput) (dto.PortalUser, error) {
	user, err := s.buildUser(ctx, common.NewID(), input, true, "")
	if err != nil {
		return dto.PortalUser{}, err
	}
	created, err := s.users.CreateUser(ctx, user)
	if err != nil {
		return dto.PortalUser{}, err
	}
	return s.mapper.ToPortalUser(created), nil
}

func (s *UserAdminService) UpdateUser(ctx context.Context, id string, input dto.CreatePortalUserInput) (dto.PortalUser, error) {
	existing, err := s.users.GetUser(ctx, id)
	if err != nil {
		return dto.PortalUser{}, err
	}
	user, err := s.buildUser(ctx, id, input, false, existing.PasswordHash)
	if err != nil {
		return dto.PortalUser{}, err
	}
	updated, err := s.users.UpdateUser(ctx, user)
	if err != nil {
		return dto.PortalUser{}, err
	}
	return s.mapper.ToPortalUser(updated), nil
}

func (s *UserAdminService) DeleteUser(ctx context.Context, id string) error {
	return s.users.DeleteUser(ctx, id)
}

func (s *UserAdminService) buildUser(ctx context.Context, id string, input dto.CreatePortalUserInput, passwordRequired bool, fallbackHash string) (model.PortalUser, error) {
	username := strings.TrimSpace(input.Username)
	if username == "" {
		return model.PortalUser{}, apperrors.NewValidation("username is required")
	}
	if passwordRequired && input.Password == "" {
		return model.PortalUser{}, apperrors.NewValidation("password is required")
	}
	for _, groupID := range input.GroupIDs {
		if _, err := s.groups.GetGroup(ctx, groupID); err != nil {
			return model.PortalUser{}, apperrors.NewValidation("one or more groupIds do not exist")
		}
	}
	passwordHash := fallbackHash
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.PortalUser{}, err
		}
		passwordHash = string(hash)
	}
	active := true
	if input.Active != nil {
		active = *input.Active
	}
	role := strings.TrimSpace(input.Role)
	if role == "" {
		role = "USER"
	}
	return model.PortalUser{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		Email:        strings.TrimSpace(input.Email),
		DisplayName:  strings.TrimSpace(input.DisplayName),
		Role:         role,
		Active:       active,
		GroupIDs:     append([]string(nil), input.GroupIDs...),
	}, nil
}
