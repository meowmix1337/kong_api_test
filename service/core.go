package service

import "github.com/jinzhu/gorm"

// ServiceCore .
type ServiceCore struct {
	DB *gorm.DB
}

// NewServiceCore .
func NewServiceCore(db *gorm.DB) *ServiceCore {
	s := new(ServiceCore)
	s.DB = db
	return s
}
