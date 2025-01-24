package routes

import (
	"github.com/Fyefhqdishka/LocFinder/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(r *mux.Router, h handlers.LocHandler) {
	LocRoutes(r, h)
}

func LocRoutes(r *mux.Router, h handlers.LocHandler) {
	r.HandleFunc("/location", h.GetLocationByIP).Methods("GET", "OPTIONS")
	r.HandleFunc("/location/{ip}", h.GetLocationForProvidedIP).Methods("GET", "OPTIONS")
	r.HandleFunc("/location/{ip}", h.UpdateLocation).Methods("PUT", "OPTIONS")
	r.HandleFunc("/location/{ip}", h.DeleteLocation).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/locations", h.GetAllLocations).Methods("GET", "OPTIONS")

	r.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}
	}).Methods("OPTIONS")
}
