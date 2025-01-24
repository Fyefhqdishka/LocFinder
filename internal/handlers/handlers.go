package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Fyefhqdishka/LocFinder/internal/models"
	"github.com/Fyefhqdishka/LocFinder/internal/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type LocHandler struct {
	Service service.ServiceInterface
	log     *slog.Logger
}

func NewLocHandler(Service service.ServiceInterface, log *slog.Logger) *LocHandler {
	return &LocHandler{Service: Service, log: log}
}

func getExternalIP() (string, error) {
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

func (h *LocHandler) GetLocationByIP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")

	if ip == "" {
		var err error
		ip, err = getExternalIP()
		if err != nil {
			h.response(w, SendError("Unable to retrieve external IP: "+err.Error()), http.StatusInternalServerError)
			return
		}
	}

	location, err := h.Service.GetLocationByIP(ip)
	if err != nil {
		h.response(w, SendError("Can't get location: "+err.Error()), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(location), http.StatusOK)
}

func (h *LocHandler) GetLocationForProvidedIP(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	location, err := h.Service.GetLocationByIP(ip)
	if err != nil {
		h.response(w, SendError("Can't get location: "+err.Error()), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(location), http.StatusOK)
}

func (h *LocHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]

	var location models.IPLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		h.response(w, SendError("Invalid request body"), http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateLocation(ip, location.Country, location.City)
	if err != nil {
		h.response(w, SendError("Can't update location: "+err.Error()), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess("Location updated"), http.StatusOK)
}

func (h *LocHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]

	err := h.Service.DeleteLocation(ip)
	if err != nil {
		h.response(w, SendError("Can't delete location: "+err.Error()), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess("Location deleted"), http.StatusOK)
}

func (h *LocHandler) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.Service.GetAllLocations()
	if err != nil {
		h.response(w, SendError("Can't fetch all locations: "+err.Error()), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(locations), http.StatusOK)
}
