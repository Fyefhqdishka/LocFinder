package repositoryInterfaces

import "github.com/Fyefhqdishka/LocFinder/internal/models"

type Storage interface {
	GetByIP(ip string) (string, string, error)
	Save(ip, country, city string) error
	Update(ip, country, city string) error
	Delete(ip string) error
	GetAll() ([]models.IPLocation, error)
}
