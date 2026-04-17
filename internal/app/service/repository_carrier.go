package service

import "plexus-bff-service-go/internal/app/repository"

type RepositoryCarrier struct {
	Users       repository.UserRepository
	Groups      repository.GroupRepository
	Permissions repository.PermissionRepository
	Domains     repository.DomainRepository
}
