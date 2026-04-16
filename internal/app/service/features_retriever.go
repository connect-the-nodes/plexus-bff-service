package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	"gopkg.in/yaml.v3"

	appcfg "plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/model"
)

type featuresRetriever struct {
	cfg          *appcfg.Config
	mapper       *mapper.FeatureMapper
	cached       []model.FeatureFlag
	cacheMutex   sync.RWMutex
}

func NewFeaturesRetriever(cfg *appcfg.Config, featureMapper *mapper.FeatureMapper) (FeaturesRetriever, error) {
	return &featuresRetriever{
		cfg:    cfg,
		mapper: featureMapper,
	}, nil
}

func (r *featuresRetriever) RetrieveFeatures(ctx context.Context) ([]model.FeatureFlag, error) {
	if r.cfg.Features.Session.Enabled {
		r.cacheMutex.RLock()
		if r.cached != nil {
			defer r.cacheMutex.RUnlock()
			return append([]model.FeatureFlag(nil), r.cached...), nil
		}
		r.cacheMutex.RUnlock()
	}

	var (
		features []model.FeatureFlag
		err error
	)
	if r.cfg.AWS.AppConfig.Features.Enabled {
		features, err = r.retrieveFromAppConfig(ctx)
	} else {
		features, err = r.retrieveFromFile()
	}
	if err != nil {
		return nil, err
	}

	if r.cfg.Features.Session.Enabled {
		r.cacheMutex.Lock()
		r.cached = append([]model.FeatureFlag(nil), features...)
		r.cacheMutex.Unlock()
	}
	return features, nil
}

func (r *featuresRetriever) IsActive(ctx context.Context, name string) (bool, error) {
	features, err := r.RetrieveFeatures(ctx)
	if err != nil {
		return false, err
	}
	for _, feature := range features {
		if feature.Name == name {
			return feature.Enabled, nil
		}
	}
	return false, nil
}

func (r *featuresRetriever) retrieveFromFile() ([]model.FeatureFlag, error) {
	filename := r.cfg.Features.File
	if filename == "" {
		filename = "features.yml"
	}
	content, err := os.ReadFile(filepath.Join("configs", filename))
	if err != nil {
		return nil, err
	}
	var payload dto.FeatureList
	if err := yaml.Unmarshal(content, &payload); err != nil {
		return nil, err
	}
	return r.mapper.MapFromDTO(payload.Features), nil
}

func (r *featuresRetriever) retrieveFromAppConfig(ctx context.Context) ([]model.FeatureFlag, error) {
	if r.cfg.AWS.AppConfig.Features.ApplicationID == "" ||
		r.cfg.AWS.AppConfig.Features.EnvironmentID == "" ||
		r.cfg.AWS.AppConfig.Features.ConfigurationID == "" {
		return nil, errors.New("aws appconfig identifiers are not configured for feature flags")
	}

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := appconfigdata.NewFromConfig(awsCfg)
	session, err := client.StartConfigurationSession(ctx, &appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:           &r.cfg.AWS.AppConfig.Features.ApplicationID,
		EnvironmentIdentifier:           &r.cfg.AWS.AppConfig.Features.EnvironmentID,
		ConfigurationProfileIdentifier:  &r.cfg.AWS.AppConfig.Features.ConfigurationID,
	})
	if err != nil {
		return nil, err
	}
	latest, err := client.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: session.InitialConfigurationToken,
	})
	if err != nil {
		return nil, err
	}

	var payload dto.FeatureList
	if err := yaml.Unmarshal(latest.Configuration, &payload); err != nil {
		return nil, err
	}
	return r.mapper.MapFromDTO(payload.Features), nil
}
