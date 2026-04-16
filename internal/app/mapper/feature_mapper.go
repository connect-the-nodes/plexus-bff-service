package mapper

import (
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/model"
)

type FeatureMapper struct{}

func NewFeatureMapper() *FeatureMapper {
	return &FeatureMapper{}
}

func (m *FeatureMapper) MapFromDTO(features []dto.FeatureFlag) []model.FeatureFlag {
	result := make([]model.FeatureFlag, 0, len(features))
	for _, feature := range features {
		parent := feature.Parent
		if parent == "" {
			parent = feature.Name
		}
		result = append(result, model.FeatureFlag{
			Name:        feature.Name,
			Parent:      parent,
			Enabled:     feature.Enabled,
			Description: feature.Description,
		})
	}
	return result
}
