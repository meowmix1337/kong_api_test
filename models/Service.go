package models

import (
	"strings"
	"time"

	"kong_api/endpoint"

	"database/sql"

	"github.com/jinzhu/gorm"
)

// Service .
type Service struct {
	ID          uint         `gorm:"primary_key;auto_increment"`
	Name        string       `gorm:"size:255;not null"`
	Description string       `gorm:"type:text"`
	Versions    []Version    `gorm:"ForeignKey:ServiceID"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time    `gorm:"default:null"`
	DeletedAt   sql.NullTime `gorm:"default:null"`
}

// ServiceAugmented .
type ServiceAugmented struct {
	ID           uint
	Name         string
	Description  string
	VersionCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

// ServiceVersions .
type ServiceVersions struct {
	ID          uint
	Name        string
	Description string
	Versions    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

func (s *ServiceVersions) serviceVersionsToEndpoint() *endpoint.ServiceVersions {
	eServiceVersions := new(endpoint.ServiceVersions)

	eServiceVersions.ID = s.ID
	eServiceVersions.Name = s.Name
	eServiceVersions.Description = s.Description
	eServiceVersions.Versions = strings.Split(s.Versions, ",")
	eServiceVersions.CreatedAt = s.CreatedAt
	eServiceVersions.UpdatedAt = s.UpdatedAt

	return eServiceVersions
}

func (s *ServiceAugmented) serviceAugmentedToEndpoint() *endpoint.ServiceIndex {
	eService := new(endpoint.ServiceIndex)

	eService.ID = s.ID
	eService.Name = s.Name
	eService.Description = s.Description
	eService.VersionCount = s.VersionCount
	eService.CreatedAt = s.CreatedAt
	eService.UpdatedAt = s.UpdatedAt

	if s.DeletedAt.Valid {
		eService.DeletedAt = &s.DeletedAt.Time
	}

	return eService
}

func (s *Service) toEndpoint() *endpoint.Service {
	eService := new(endpoint.Service)

	eService.ID = s.ID
	eService.Name = s.Name
	eService.Description = s.Description
	eService.CreatedAt = s.CreatedAt
	eService.UpdatedAt = s.UpdatedAt

	if s.DeletedAt.Valid {
		eService.DeletedAt = &s.DeletedAt.Time
	}

	return eService
}

// ExistsByName .
func (s *Service) ExistsByName(db *gorm.DB) (bool, error) {
	var err error
	err = db.Debug().Model(Service{}).Where("name = ? AND deleted_at IS NULL", s.Name).Take(&s).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return s != nil, nil
}

// ExistsByID .
func (s *Service) ExistsByID(db *gorm.DB) (bool, error) {
	var err error
	err = db.Debug().Model(Service{}).Where("id = ? AND deleted_at IS NULL", s.ID).Take(&s).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return s != nil, nil
}

// Store .
func (s *Service) Store(db *gorm.DB) (*endpoint.Service, error) {

	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return nil, err
	}
	return s.toEndpoint(), nil
}

// AllFilter .
func (s *Service) CountAll(db *gorm.DB) (int, error) {
	var err error
	var count int

	err = db.Debug().Model(&Service{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}

// AllFilter .
func (s *Service) AllFilter(db *gorm.DB, keyword, sortBy, sortOrder string, offset, pageSize int) ([]endpoint.ServiceIndex, error) {
	var err error
	var services []ServiceAugmented

	chain := db.Debug().Table("services").Select("services.id, services.name, services.description, services.created_at, services.updated_at, services.deleted_at, COUNT(versions.id) as version_count").
		Joins("JOIN versions ON versions.service_id = services.id").Group("services.id")

	if keyword != "" {
		chain = chain.Where("name LIKE ?", "%"+keyword+"%").
			Or("description LIKE ?", "%"+keyword+"%")
	}

	if offset > 0 {
		chain = chain.Offset(offset)
	}

	if pageSize > 0 {
		chain = chain.Limit(pageSize)
	}

	if sortBy != "" {
		chain = chain.Order(sortBy + " " + sortOrder)
	}

	err = chain.Find(&services).Error

	if err != nil {
		return nil, err
	}

	var eServices []endpoint.ServiceIndex
	for _, dbService := range services {
		eServices = append(eServices, *dbService.serviceAugmentedToEndpoint())
	}

	return eServices, err
}

// ByID .
func (s *Service) ByID(db *gorm.DB, id uint32) (*endpoint.ServiceVersions, error) {
	var err error
	var serviceVersions ServiceVersions

	err = db.Debug().Table("services").Select("services.id, services.name, services.description, services.created_at, services.updated_at, services.deleted_at, GROUP_CONCAT(versions.version) as versions").
		Joins("JOIN versions ON versions.service_id = services.id").Where("services.id = ?", id).Group("services.id").Take(&serviceVersions).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, endpoint.ErrServiceNotFound
		}
		return nil, err
	}

	return serviceVersions.serviceVersionsToEndpoint(), err
}

// Update .
func (s *Service) Update(db *gorm.DB, uid uint32) (*endpoint.Service, error) {

	db = db.Debug().Model(&Service{}).Where("id = ?", uid).Take(&Service{}).UpdateColumns(
		map[string]interface{}{
			"name":        s.Name,
			"description": s.Description,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return nil, db.Error
	}

	// This is the display the updated service
	err := db.Debug().Model(&Service{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return nil, err
	}

	return s.toEndpoint(), nil
}

// DeleteByID .
func (s *Service) DeleteByID(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Service{}).Where("id = ?", uid).Take(&Service{}).UpdateColumns(
		map[string]interface{}{
			"deleted_at": time.Now(),
		},
	)

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
