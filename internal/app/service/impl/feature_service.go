package impl

import (
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/service"
)

type FeatureService struct{}

func NewFeatureService() service.FeatureService {
	return &FeatureService{}
}

func (s *FeatureService) GetStatus() dto.FeatureResponse {
	return dto.FeatureResponse{
		Status:  "OK",
		Message: "Service is running",
	}
}
