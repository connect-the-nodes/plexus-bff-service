package service

import (
	"context"

	"plexus-bff-service-go/internal/app/dto"
)

type FeatureFlagsService struct {
	retriever FeaturesRetriever
}

func NewFeatureFlagsService(retriever FeaturesRetriever) *FeatureFlagsService {
	return &FeatureFlagsService{retriever: retriever}
}

func (s *FeatureFlagsService) RetrieveFeatures(ctx context.Context) (dto.FeatureFlagsResponse, error) {
	features, err := s.retriever.RetrieveFeatures(ctx)
	if err != nil {
		return dto.FeatureFlagsResponse{}, err
	}

	result := make([]dto.FeatureFlag, 0, len(features))
	for _, feature := range features {
		result = append(result, dto.FeatureFlag{
			Name:        feature.Name,
			Parent:      feature.Parent,
			Enabled:     feature.Enabled,
			Description: feature.Description,
		})
	}
	return dto.FeatureFlagsResponse{Features: result}, nil
}
