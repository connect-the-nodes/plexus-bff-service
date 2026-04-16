package app

import (
	"net/http"

	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/controller"
	"plexus-bff-service-go/internal/app/httpx"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/repository/impl"
	"plexus-bff-service-go/internal/app/service"
	serviceimpl "plexus-bff-service-go/internal/app/service/impl"
)

type Application struct {
	handler http.Handler
}

func New(cfg *config.Config) (*Application, error) {
	featureMapper := mapper.NewFeatureMapper()
	featuresRetriever, err := service.NewFeaturesRetriever(cfg, featureMapper)
	if err != nil {
		return nil, err
	}

	connectorsRepository := impl.NewConnectorsRepository()
	connectorsService := serviceimpl.NewConnectorsService(connectorsRepository, featuresRetriever)
	featureService := serviceimpl.NewFeatureService()
	featureFlagsService := service.NewFeatureFlagsService(featuresRetriever)
	observabilityService := serviceimpl.NewObservabilityService(cfg.Observability.Service.BaseURL)
	authController := controller.NewAuthController(cfg.Auth.Cognito)
	connectorsController := controller.NewConnectorsController(connectorsService)
	featureController := controller.NewFeatureController(featureService)
	featureFlagsController := controller.NewFeatureFlagsController(featureFlagsService)
	observabilityController := controller.NewObservabilityController(observabilityService)
	adminObservabilityController := controller.NewAdminObservabilityController(observabilityService)
	testSessionController := controller.NewTestSessionController()

	router, err := httpx.NewRouter(
		cfg,
		authController,
		connectorsController,
		featureController,
		featureFlagsController,
		observabilityController,
		adminObservabilityController,
		testSessionController,
	)
	if err != nil {
		return nil, err
	}

	return &Application{handler: router}, nil
}

func (a *Application) Handler() http.Handler {
	return a.handler
}
