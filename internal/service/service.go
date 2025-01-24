package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Fyefhqdishka/LocFinder/internal/models"
	"github.com/Fyefhqdishka/LocFinder/internal/storage/repositoryInterfaces"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type ServiceInterface interface {
	GetLocationByIP(ip string) (*models.IPLocation, error)
	UpdateLocation(ip, country, city string) error
	DeleteLocation(ip string) error
	GetAllLocations() ([]models.IPLocation, error)
	GetExternalIP() (string, error)
	fetchFromAPI(ip string) (models.IPLocation, error)
}

type LocService struct {
	repo repositoryInterfaces.Storage
	log  *slog.Logger
}

func NewLocService(repo repositoryInterfaces.Storage, log *slog.Logger) *LocService {
	return &LocService{repo: repo, log: log}
}

func (s *LocService) GetExternalIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("could not get external IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %v", err)
	}

	return string(body), nil
}

func (s *LocService) GetLocationByIP(ip string) (*models.IPLocation, error) {
	country, city, err := s.repo.GetByIP(ip)
	if err == nil {
		return &models.IPLocation{IP: ip, Country: country, City: city}, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("Ошибка при поиске локации по IP в базе", "error", err)
		return nil, err
	}

	location, err := s.fetchFromAPI(ip)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(location.IP, location.Country, location.City)
	if err != nil {
		return nil, err
	}

	s.log.Debug("Локация сохранена в базе", "ip", ip, "location", location)

	return &location, nil
}

func (s *LocService) UpdateLocation(ip, country, city string) error {
	err := s.repo.Update(ip, country, city)
	if err != nil {
		s.log.Error("Ошибка при обновлении локации", "error", err)
		return err
	}
	s.log.Debug("Локация обновлена в базе", "ip", ip)
	return nil
}

func (s *LocService) DeleteLocation(ip string) error {
	err := s.repo.Delete(ip)
	if err != nil {
		s.log.Error("Ошибка при удалении локации", "error", err)
		return err
	}
	s.log.Debug("Локация удалена из базы", "ip", ip)
	return nil
}

func (s *LocService) GetAllLocations() ([]models.IPLocation, error) {
	locations, err := s.repo.GetAll()
	if err != nil {
		s.log.Error("Ошибка при получении всех локаций", "error", err)
		return nil, err
	}
	s.log.Debug("Получены все локации", "count", len(locations))
	return locations, nil
}

func (s *LocService) fetchFromAPI(ip string) (models.IPLocation, error) {
	apiURL := "http://ip-api.com/json/" + ip
	resp, err := http.Get(apiURL)
	if err != nil {
		return models.IPLocation{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.IPLocation{}, err
	}

	s.log.Debug("Ответ от API", "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return models.IPLocation{}, errors.New("некорректный ответ от API: " + resp.Status)
	}

	var location models.IPLocation
	if err := json.Unmarshal(body, &location); err != nil {
		return models.IPLocation{}, err
	}

	return location, nil
}
