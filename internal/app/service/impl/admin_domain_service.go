package impl

import (
	"context"
	"strings"
	"time"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/common"
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
	"plexus-bff-service-go/internal/app/service"
)

const (
	domainStatusDraft         = "DRAFT"
	domainStatusPendingReview = "PENDING_APPROVAL"
	domainStatusApproved      = "APPROVED"
	domainStatusRejected      = "REJECTED"
)

type DomainAdminService struct {
	domains repository.DomainRepository
	groups  repository.GroupRepository
	mapper  *mapper.AdminMapper
}

func NewDomainAdminService(domains repository.DomainRepository, groups repository.GroupRepository) service.DomainAdminService {
	return &DomainAdminService{domains: domains, groups: groups, mapper: mapper.NewAdminMapper()}
}

func (s *DomainAdminService) GetWorkspace(ctx context.Context) (dto.DomainWorkspace, error) {
	groups, err := s.groups.ListGroups(ctx)
	if err != nil {
		return dto.DomainWorkspace{}, err
	}
	domains, err := s.domains.ListDomains(ctx)
	if err != nil {
		return dto.DomainWorkspace{}, err
	}
	return dto.DomainWorkspace{
		OwnerGroups: s.mapper.ToPortalGroups(groups),
		Domains:     s.mapper.ToRegisteredDomains(domains),
	}, nil
}

func (s *DomainAdminService) ListDomains(ctx context.Context) ([]dto.RegisteredDomain, error) {
	items, err := s.domains.ListDomains(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToRegisteredDomains(items), nil
}

func (s *DomainAdminService) ListApprovedDomains(ctx context.Context) ([]dto.RegisteredDomain, error) {
	items, err := s.domains.ListApprovedDomains(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToRegisteredDomains(items), nil
}

func (s *DomainAdminService) GetDomain(ctx context.Context, id string) (dto.RegisteredDomain, error) {
	item, err := s.domains.GetDomain(ctx, id)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	return s.mapper.ToRegisteredDomain(item), nil
}

func (s *DomainAdminService) CreateDomain(ctx context.Context, input dto.DomainMutationRequest) (dto.RegisteredDomain, error) {
	domain, err := s.buildDomain(ctx, common.NewID(), input)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	item, err := s.domains.CreateDomain(ctx, domain)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	return s.mapper.ToRegisteredDomain(item), nil
}

func (s *DomainAdminService) UpdateDomain(ctx context.Context, id string, input dto.DomainMutationRequest) (dto.RegisteredDomain, error) {
	existing, err := s.domains.GetDomain(ctx, id)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	if existing.Status == domainStatusApproved && normalizeMode(input.Mode) != domainStatusDraft {
		return dto.RegisteredDomain{}, apperrors.NewValidation("approved domains must be edited as draft before review")
	}
	domain, err := s.buildDomain(ctx, id, input)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	domain.Review = existing.Review
	item, err := s.domains.UpdateDomain(ctx, domain)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	return s.mapper.ToRegisteredDomain(item), nil
}

func (s *DomainAdminService) DeleteDomain(ctx context.Context, id string) error {
	item, err := s.domains.GetDomain(ctx, id)
	if err != nil {
		return err
	}
	if item.Status != domainStatusDraft {
		return apperrors.NewValidation("only draft domains can be deleted")
	}
	return s.domains.DeleteDomain(ctx, id)
}

func (s *DomainAdminService) ReviewDomain(ctx context.Context, id string, input dto.ReviewDecisionRequest) (dto.RegisteredDomain, error) {
	if strings.TrimSpace(input.Author) == "" || strings.TrimSpace(input.Comment) == "" {
		return dto.RegisteredDomain{}, apperrors.NewValidation("author and comment are required")
	}
	domain, err := s.domains.GetDomain(ctx, id)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	if domain.Status != domainStatusPendingReview {
		return dto.RegisteredDomain{}, apperrors.NewValidation("only pending review domains can be reviewed")
	}
	decision := strings.ToUpper(strings.TrimSpace(input.Decision))
	switch decision {
	case domainStatusApproved, domainStatusRejected:
	default:
		return dto.RegisteredDomain{}, apperrors.NewValidation("decision must be APPROVED or REJECTED")
	}
	item, err := s.domains.AddDomainReview(ctx, id, model.ReviewComment{
		ID:        common.NewID(),
		Author:    strings.TrimSpace(input.Author),
		AuthorID:  strings.TrimSpace(input.AuthorID),
		Decision:  decision,
		Comment:   strings.TrimSpace(input.Comment),
		CreatedAt: time.Now().UTC(),
	}, decision)
	if err != nil {
		return dto.RegisteredDomain{}, err
	}
	return s.mapper.ToRegisteredDomain(item), nil
}

func (s *DomainAdminService) buildDomain(ctx context.Context, id string, input dto.DomainMutationRequest) (model.RegisteredDomain, error) {
	name := strings.TrimSpace(input.Domain.Name)
	if name == "" {
		return model.RegisteredDomain{}, apperrors.NewValidation("domain.name is required")
	}
	if strings.TrimSpace(input.Domain.OwnerGroupID) == "" {
		return model.RegisteredDomain{}, apperrors.NewValidation("domain.ownerGroupId is required")
	}
	if _, err := s.groups.GetGroup(ctx, input.Domain.OwnerGroupID); err != nil {
		return model.RegisteredDomain{}, apperrors.NewValidation("ownerGroupId does not exist")
	}
	return model.RegisteredDomain{
		ID:           id,
		Name:         name,
		Description:  strings.TrimSpace(input.Domain.Description),
		OwnerGroupID: input.Domain.OwnerGroupID,
		Status:       normalizeMode(input.Mode),
		Metadata:     cloneDomainMetadata(input.Domain.Metadata),
		Review:       []model.ReviewComment{},
	}, nil
}

func normalizeMode(mode string) string {
	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case "", domainStatusDraft:
		return domainStatusDraft
	case "SUBMIT":
		return domainStatusPendingReview
	default:
		return domainStatusDraft
	}
}

func cloneDomainMetadata(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	result := make(map[string]any, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}
