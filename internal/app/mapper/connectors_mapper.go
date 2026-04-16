package mapper

import (
	"plexus-bff-service-go/internal/app/dto"
	"plexus-bff-service-go/internal/app/model"
)

type ConnectorsMapper struct{}

func NewConnectorsMapper() *ConnectorsMapper {
	return &ConnectorsMapper{}
}

func (m *ConnectorsMapper) Map(info model.ConnectorsInfo) dto.ConnectorsResponse {
	return dto.ConnectorsResponse{
		Status:  info.Status,
		Message: info.Message,
	}
}
