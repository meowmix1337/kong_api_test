package service

import (
	"kong_api/endpoint"
	"kong_api/models"
)

// IServicesService .
type IServicesService interface {
	All(keyword, sortBy, sortOrder string, page, pageSize int) ([]endpoint.ServiceIndex, int, error)
	ByID(id uint32) (*endpoint.ServiceVersions, error)
	Store(newService endpoint.Service) (*endpoint.Service, error)
	Update(id uint32, updatedService endpoint.Service) (*endpoint.Service, error)
	DeleteByID(id uint32) error
}

// ServicesService .
type ServicesService struct {
	*ServiceCore
}

// NewServicesService .
func NewServicesService(sCore *ServiceCore) *ServicesService {
	s := new(ServicesService)

	s.ServiceCore = sCore

	return s
}

// All .
func (s *ServicesService) All(keyword, sortBy, sortOrder string, page, pageSize int) ([]endpoint.ServiceIndex, int, error) {
	dbService := models.Service{}

	if pageSize < 1 {
		// default to page per page
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	// default sorts
	if sortBy == "" {
		sortBy = "id"
	}

	if sortOrder == "" {
		sortOrder = "ASC"
	}

	// get TOTAL count
	totalServices, err := dbService.CountAll(s.DB)
	if err != nil {
		return nil, 0, err
	}

	eServices, err := dbService.AllFilter(s.DB, keyword, sortBy, sortOrder, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return eServices, totalServices, nil
}

// ByID .
func (s *ServicesService) ByID(id uint32) (*endpoint.ServiceVersions, error) {
	dbService := models.Service{}
	eServiceVersions, err := dbService.ByID(s.DB, id)
	if err != nil {
		return nil, err
	}

	return eServiceVersions, nil
}

// Store .
func (s *ServicesService) Store(newService endpoint.Service) (*endpoint.Service, error) {
	dbService := ToDBModel(newService)

	// duplicate names are allowed

	// TODO: Add transaction
	savedService, err := dbService.Store(s.DB)
	if err != nil {
		return nil, err
	}

	initialVersion := &models.Version{
		Version:   newService.InitialVersion,
		ServiceID: savedService.ID,
	}
	newVersion, err := initialVersion.Store(s.DB)
	if err != nil {
		return nil, err
	}

	savedService.InitialVersion = newVersion.Version

	return savedService, nil
}

// Update .
func (s *ServicesService) Update(id uint32, updatedService endpoint.Service) (*endpoint.Service, error) {
	dbService := ToDBModel(updatedService)

	// check if service exists
	existService := models.Service{}
	_, err := existService.ByID(s.DB, id)
	if err != nil {
		return nil, err
	}

	eService, err := dbService.Update(s.DB, uint32(existService.ID))
	if err != nil {
		return nil, err
	}

	return eService, nil
}

// DeleteByID .
func (s *ServicesService) DeleteByID(id uint32) error {

	service := models.Service{}

	// check if service exists
	eService, err := service.ByID(s.DB, id)
	if err != nil {
		return err
	}

	_, err = service.DeleteByID(s.DB, uint32(eService.ID))
	if err != nil {
		return err
	}

	return nil
}

// ToDBModel .
func ToDBModel(s endpoint.Service) *models.Service {
	dbService := new(models.Service)

	dbService.ID = s.ID
	dbService.Name = s.Name
	dbService.Description = s.Description
	dbService.CreatedAt = s.CreatedAt
	dbService.UpdatedAt = s.UpdatedAt

	return dbService
}
