package impl

import (
	"context"
	"errors"

	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/feature"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/repository"
	"plexus-bff-service-go/internal/app/service"
)

type ConnectorsService struct {
	repository repository.ConnectorsRepository
	mapper     *mapper.ConnectorsMapper
	retriever  service.FeaturesRetriever
}

func NewConnectorsService(repo repository.ConnectorsRepository, retriever service.FeaturesRetriever) service.ConnectorsService {
	return &ConnectorsService{
		repository: repo,
		mapper:     mapper.NewConnectorsMapper(),
		retriever:  retriever,
	}
}

func (s *ConnectorsService) GetConnectors(ctx context.Context) (dto.ConnectorsResponse, error) {
	active, err := s.retriever.IsActive(ctx, feature.Connectors)
	if err != nil {
		return dto.ConnectorsResponse{}, err
	}
	if !active {
		return dto.ConnectorsResponse{}, errors.New("feature FEATURE_CONNECTORS is not enabled")
	}
	return s.mapper.Map(s.repository.FetchConnectors()), nil
}
