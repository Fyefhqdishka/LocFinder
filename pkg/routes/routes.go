package routes

import (
	"github.com/Fyefhqdishka/LocFinder/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h handlers.LocHandler) {
	LocRoutes(r, h)
}

func LocRoutes(r *mux.Router, h handlers.LocHandler) {
	r.HandleFunc("/location", h.GetLocationByIP).Methods("GET")
	r.HandleFunc("/location/{ip}", h.GetLocationForProvidedIP).Methods("GET")
	r.HandleFunc("/location/{ip}", h.UpdateLocation).Methods("PUT")
	r.HandleFunc("/location/{ip}", h.DeleteLocation).Methods("DELETE")
	r.HandleFunc("/locations", h.GetAllLocations).Methods("GET")
}
