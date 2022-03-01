package service

import (
	"kong_api/endpoint"
	"kong_api/models"
)

// IVersionsService .
type IVersionsService interface {
	Store(newService endpoint.Version) (*endpoint.Version, error)
}

// VersionsService .
type VersionsService struct {
	*ServiceCore
}

// NewVerionsService .
func NewVerionsService(sCore *ServiceCore) *VersionsService {
	s := new(VersionsService)

	s.ServiceCore = sCore

	return s
}

// Store .
func (s *VersionsService) Store(newVersion endpoint.Version) (*endpoint.Version, error) {
	dbVersion := s.ToDBModel(newVersion)

	savedVersion, err := dbVersion.Store(s.DB)
	if err != nil {
		return nil, err
	}

	return savedVersion, nil
}

// ToDBModel .
func (s *VersionsService) ToDBModel(v endpoint.Version) *models.Version {
	dbVersion := new(models.Version)

	dbVersion.ID = v.ID
	dbVersion.Version = v.Version
	dbVersion.ServiceID = v.ServiceID
	dbVersion.CreatedAt = v.CreatedAt
	dbVersion.UpdatedAt = v.UpdatedAt

	return dbVersion
}
