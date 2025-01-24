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
	FetchFromAPI(ip string) (models.IPLocation, error)
}

type LocService struct {
	repo repositoryInterfaces.Storage
	log  *slog.Logger
}

func NewLocService(repo repositoryInterfaces.Storage, log *slog.Logger) *LocService {
	return &LocService{repo: repo, log: log}
}

func (s *LocService) GetExternalIP() (string, error) {
	s.log.Debug("Попытка получить внешний IP через API", "url", "https://api.ipify.org")
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		s.log.Error("Не удалось получить внешний IP", "error", err)
		return "", fmt.Errorf("не удалось получить внешний IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Не удалось прочитать тело ответа", "error", err)
		return "", fmt.Errorf("не удалось прочитать тело ответа: %v", err)
	}

	s.log.Debug("Успешно получен внешний IP", "ip", string(body))
	return string(body), nil
}

func (s *LocService) GetLocationByIP(ip string) (*models.IPLocation, error) {
	s.log.Debug("Поиск локации для IP", "ip", ip)

	country, city, err := s.repo.GetByIP(ip)
	if err == nil {
		s.log.Debug("Локация найдена в базе данных", "ip", ip, "country", country, "city", city)
		return &models.IPLocation{IP: ip, Country: country, City: city}, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("Ошибка при поиске локации по IP в базе данных", "ip", ip, "error", err)
		return nil, err
	}

	s.log.Debug("Локация не найдена в базе данных, пытаемся получить с внешнего API", "ip", ip)
	location, err := s.FetchFromAPI(ip)
	if err != nil {
		s.log.Error("Не удалось получить локацию с API", "ip", ip, "error", err)
		return nil, err
	}

	err = s.repo.Save(location.IP, location.Country, location.City)
	if err != nil {
		s.log.Error("Не удалось сохранить локацию в базе данных", "ip", ip, "error", err)
		return nil, err
	}

	s.log.Debug("Локация успешно сохранена в базе данных", "ip", ip, "country", location.Country, "city", location.City)
	return &location, nil
}

func (s *LocService) UpdateLocation(ip, country, city string) error {
	s.log.Debug("Обновление локации", "ip", ip, "country", country, "city", city)
	err := s.repo.Update(ip, country, city)
	if err != nil {
		s.log.Error("Ошибка при обновлении локации", "error", err)
		return err
	}
	s.log.Debug("Локация обновлена в базе данных", "ip", ip)
	return nil
}

func (s *LocService) DeleteLocation(ip string) error {
	s.log.Debug("Удаление локации", "ip", ip)
	err := s.repo.Delete(ip)
	if err != nil {
		s.log.Error("Ошибка при удалении локации", "error", err)
		return err
	}
	s.log.Debug("Локация удалена из базы данных", "ip", ip)
	return nil
}

func (s *LocService) GetAllLocations() ([]models.IPLocation, error) {
	s.log.Debug("Получение всех локаций из базы данных")
	locations, err := s.repo.GetAll()
	if err != nil {
		s.log.Error("Ошибка при получении всех локаций", "error", err)
		return nil, err
	}
	s.log.Debug("Получены все локации из базы данных", "count", len(locations))
	return locations, nil
}

func (s *LocService) FetchFromAPI(ip string) (models.IPLocation, error) {
	s.log.Debug("Запрос локации с внешнего API", "ip", ip)
	apiURL := "http://ip-api.com/json/" + ip
	resp, err := http.Get(apiURL)
	if err != nil {
		s.log.Error("Ошибка при запросе к API", "error", err)
		return models.IPLocation{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Не удалось прочитать тело ответа от API", "error", err)
		return models.IPLocation{}, err
	}

	s.log.Debug("Ответ от API получен", "body", string(body))

	if resp.StatusCode != http.StatusOK {
		s.log.Error("Некорректный ответ от API", "status", resp.Status)
		return models.IPLocation{}, errors.New("некорректный ответ от API: " + resp.Status)
	}

	var location models.IPLocation
	if err := json.Unmarshal(body, &location); err != nil {
		s.log.Error("Ошибка при разборе ответа API", "error", err)
		return models.IPLocation{}, err
	}

	s.log.Debug("Локация получена с API", "ip", location.IP, "country", location.Country, "city", location.City)
	return location, nil
}
