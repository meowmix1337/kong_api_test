package models

import (
	"kong_api/endpoint"
	"time"

	"database/sql"

	"github.com/jinzhu/gorm"
)

// Version .
type Version struct {
	ID        uint         `gorm:"primary_key;auto_increment"`
	Version   string       `gorm:"size:255;not null"`
	ServiceID uint         `gorm:"column:service_id"`
	CreatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time    `gorm:"default:null"`
	DeletedAt sql.NullTime `gorm:"default:null"`
}

func (s *Version) toEndpoint() *endpoint.Version {
	eVersion := new(endpoint.Version)

	eVersion.ID = s.ID
	eVersion.Version = s.Version
	eVersion.ServiceID = s.ServiceID
	eVersion.CreatedAt = s.CreatedAt
	eVersion.UpdatedAt = s.UpdatedAt

	if s.DeletedAt.Valid {
		eVersion.DeletedAt = &s.DeletedAt.Time
	}

	return eVersion
}

// Store .
func (v *Version) Store(db *gorm.DB) (*endpoint.Version, error) {

	var err error
	err = db.Debug().Create(&v).Error
	if err != nil {
		return nil, err
	}
	return v.toEndpoint(), nil
}
