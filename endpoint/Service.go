package endpoint

import (
	"errors"
	"time"
)

var (
	// ErrServiceAlreadyExists .
	ErrServiceAlreadyExists = errors.New("Service already exists")
	// ErrServiceNotFound .
	ErrServiceNotFound = errors.New("Service Not Found")
	// ErrNameRequired .
	ErrNameRequired = errors.New("Name Required")
	// ErrVersionRequired .
	ErrVersionRequired = errors.New("Version Required")
)

// Service represents a service with an initial version
type Service struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	InitialVersion string     `json:"initial_version"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// ServiceIndex represents a service with a count of versions
type ServiceIndex struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	VersionCount int        `json:"versions"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// ServiceVersions represents the service with list of versions
type ServiceVersions struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Versions    []string   `json:"versions"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// PaginatedServices .
type PaginatedServices struct {
	Total    int            `json:"total"`
	Services []ServiceIndex `json:"services"`
}

// Validate .
func (s *Service) Validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.InitialVersion == "" {
		return ErrVersionRequired
	}
	return nil
}
