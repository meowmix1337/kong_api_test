package endpoint

import (
	"errors"
	"time"
)

var (
	// ErrServiceIDRequired .
	ErrServiceIDRequired = errors.New("Service ID Required")
)

// Version .
type Version struct {
	ID        uint       `json:"id"`
	Version   string     `json:"version"`
	ServiceID uint       `json:"service_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Validate .
func (s *Version) Validate() error {
	if s.Version == "" {
		return ErrVersionRequired
	}
	if s.ServiceID == 0 {
		return ErrServiceIDRequired
	}
	return nil
}
