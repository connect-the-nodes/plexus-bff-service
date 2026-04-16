package impl

import (
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
)

type ConnectorsRepository struct{}

func NewConnectorsRepository() repository.ConnectorsRepository {
	return &ConnectorsRepository{}
}

func (r *ConnectorsRepository) FetchConnectors() model.ConnectorsInfo {
	return model.ConnectorsInfo{
		Status:  "OK",
		Message: "Connectors Service is running",
	}
}
