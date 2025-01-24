package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Fyefhqdishka/LocFinder/internal/handlers"
	"github.com/Fyefhqdishka/LocFinder/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Мок сервиса
type MockService struct {
	mock.Mock
}

func (m *MockService) GetExternalIP() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockService) GetLocationByIP(ip string) (*models.IPLocation, error) {
	args := m.Called(ip)
	return args.Get(0).(*models.IPLocation), args.Error(1)
}

func (m *MockService) UpdateLocation(ip, country, city string) error {
	args := m.Called(ip, country, city)
	return args.Error(0)
}

func (m *MockService) DeleteLocation(ip string) error {
	args := m.Called(ip)
	return args.Error(0)
}

func (m *MockService) GetAllLocations() ([]models.IPLocation, error) {
	args := m.Called()
	return args.Get(0).([]models.IPLocation), args.Error(1)
}

// Добавляем метод FetchFromAPI
func (m *MockService) FetchFromAPI(ip string) (models.IPLocation, error) {
	args := m.Called(ip)
	return args.Get(0).(models.IPLocation), args.Error(1)
}

func TestDeleteLocation(t *testing.T) {
	mockService := new(MockService)
	log := slog.Logger{}

	handler := handlers.NewLocHandler(mockService, &log)

	ip := "37.99.42.212"

	mockService.On("DeleteLocation", ip).Return(nil)

	req, err := http.NewRequest("DELETE", "/locations/"+ip, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/locations/{ip}", handler.DeleteLocation).Methods("DELETE")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateLocation(t *testing.T) {
	mockService := new(MockService)
	log := slog.Logger{}

	handler := handlers.NewLocHandler(mockService, &log)

	mockService.On("UpdateLocation", "37.99.42.212", "Updated Country", "Updated City").Return(nil)

	location := models.IPLocation{
		IP:      "37.99.42.212",
		Country: "Updated Country",
		City:    "Updated City",
	}
	locationJSON, err := json.Marshal(location)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/locations", bytes.NewBuffer(locationJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/locations", handler.UpdateLocation).Methods("PUT")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}

func TestGetAllLocations(t *testing.T) {
	mockService := new(MockService)
	log := slog.Logger{}

	handler := handlers.NewLocHandler(mockService, &log)

	mockService.On("GetAllLocations").Return([]models.IPLocation{
		{IP: "37.99.42.212", Country: "Kazakhstan", City: "Almaty"},
	}, nil)

	req, err := http.NewRequest("GET", "/locations", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/locations", handler.GetAllLocations).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	result, ok := response["result"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'result' to be an array, got: %v", response["result"])
	}

	location := result[0].(map[string]interface{})
	assert.Equal(t, "Kazakhstan", location["country"])
	assert.Equal(t, "Almaty", location["city"])
	assert.Equal(t, "37.99.42.212", location["query"])

	mockService.AssertExpectations(t)
}

func TestGetLocationByIP(t *testing.T) {
	mockService := new(MockService)
	log := slog.Logger{}

	handler := handlers.NewLocHandler(mockService, &log)

	mockService.On("GetExternalIP").Return("37.99.42.212", nil)

	mockService.On("GetLocationByIP", "37.99.42.212").Return(&models.IPLocation{
		IP:      "37.99.42.212",
		Country: "Kazakhstan",
		City:    "Almaty",
	}, nil)

	req, err := http.NewRequest("GET", "/locations/37.99.42.212", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/locations/{ip}", handler.GetLocationByIP).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	location := response["result"].(map[string]interface{})
	assert.Equal(t, "Kazakhstan", location["country"])
	assert.Equal(t, "Almaty", location["city"])
	assert.Equal(t, "37.99.42.212", location["query"])

	mockService.AssertExpectations(t)
}
