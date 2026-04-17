package app

import (
	"context"
	"net/http"

	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/controller"
	"plexus-bff-service-go/internal/app/database"
	"plexus-bff-service-go/internal/app/httpx"
	"plexus-bff-service-go/internal/app/mapper"
	"plexus-bff-service-go/internal/app/repository/impl"
	"plexus-bff-service-go/internal/app/service"
	serviceimpl "plexus-bff-service-go/internal/app/service/impl"
)

type Application struct {
	handler http.Handler
	closeFn func()
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

	postgres, err := database.Open(context.Background(), cfg.Database)
	if err != nil {
		return nil, err
	}
	if cfg.Database.Enabled() && cfg.Database.AutoMigrate {
		if err := database.RunMigrations(context.Background(), postgres); err != nil {
			if postgres != nil {
				postgres.Close()
			}
			return nil, err
		}
	}

	var repositories service.RepositoryCarrier
	if postgres != nil {
		adminRepository := impl.NewPostgresAdminRepository(postgres.Pool())
		repositories = service.RepositoryCarrier{
			Users:       adminRepository.Users(),
			Groups:      adminRepository.Groups(),
			Permissions: adminRepository.Permissions(),
			Domains:     adminRepository.Domains(),
		}
	} else {
		adminRepository := impl.NewMemoryAdminRepository()
		repositories = service.RepositoryCarrier{
			Users:       adminRepository.Users(),
			Groups:      adminRepository.Groups(),
			Permissions: adminRepository.Permissions(),
			Domains:     adminRepository.Domains(),
		}
	}

	userAdminService := serviceimpl.NewUserAdminService(repositories.Users, repositories.Groups)
	groupAdminService := serviceimpl.NewGroupAdminService(repositories.Groups, repositories.Permissions)
	permissionAdminService := serviceimpl.NewPermissionAdminService(repositories.Permissions)
	domainAdminService := serviceimpl.NewDomainAdminService(repositories.Domains, repositories.Groups)
	adminUsersController := controller.NewAdminUsersController(userAdminService)
	adminGroupsController := controller.NewAdminGroupsController(groupAdminService)
	adminPermissionsController := controller.NewAdminPermissionsController(permissionAdminService)
	domainsController := controller.NewDomainsController(domainAdminService)

	router, err := httpx.NewRouter(
		cfg,
		authController,
		connectorsController,
		featureController,
		featureFlagsController,
		observabilityController,
		adminObservabilityController,
		testSessionController,
		adminUsersController,
		adminGroupsController,
		adminPermissionsController,
		domainsController,
	)
	if err != nil {
		if postgres != nil {
			postgres.Close()
		}
		return nil, err
	}

	return &Application{
		handler: router,
		closeFn: func() {
			if postgres != nil {
				postgres.Close()
			}
		},
	}, nil
}

func (a *Application) Handler() http.Handler {
	return a.handler
}

func (a *Application) Close() {
	if a == nil || a.closeFn == nil {
		return
	}
	a.closeFn()
}
