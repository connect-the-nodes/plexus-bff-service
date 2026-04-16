package repository

import "plexus-bff-service-go/internal/app/model"

type ConnectorsRepository interface {
	FetchConnectors() model.ConnectorsInfo
}
