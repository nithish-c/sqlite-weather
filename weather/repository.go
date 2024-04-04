package weather

import (
	"errors"
)

var (
	ErrDuplicate = errors.New("weather: duplicate entry")
	ErrNotFound  = errors.New("weather: not found")
	ErrUpdate    = errors.New("weather: update failed")
	ErrDelete    = errors.New("weather: delete failed")
)

type Repository interface {
	Migrate() error
	Create(weather Weather) (*Weather, error)
	All() ([]Weather, error)
	GetByID(id int64) (*Weather, error)
	Update(weather Weather) (*Weather, error)
	Delete(id int64) error

	GetByCity(city string) (*Weather, error)
	DropTable() error
	UniqueCities(r *RepositorySQLite) ([]string, error)
}
