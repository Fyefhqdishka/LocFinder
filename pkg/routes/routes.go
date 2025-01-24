package routes

import (
	"github.com/Fyefhqdishka/LocFinder/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func RegisterRoutes(r *mux.Router, h handlers.LocHandler) http.Handler {
	LocRoutes(r, h)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	return corsHandler
}

func LocRoutes(r *mux.Router, h handlers.LocHandler) {
	r.HandleFunc("/location", h.GetLocationByIP).Methods("GET")
	r.HandleFunc("/location/{ip}", h.GetLocationForProvidedIP).Methods("GET")
	r.HandleFunc("/location/{ip}", h.UpdateLocation).Methods("PUT")
	r.HandleFunc("/location/{ip}", h.DeleteLocation).Methods("DELETE")
	r.HandleFunc("/locations", h.GetAllLocations).Methods("GET")
}
