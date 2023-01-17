package handler

import (
	"net/http"

	"github.com/MrDavudov/AirtableToDraw/pkg/service"
	"github.com/gorilla/mux"
)

type server struct {
	router 		*mux.Router
	services	*service.Service
}

func New(services *service.Service) *server {
	s := &server{
		router:  	mux.NewRouter(),
		services:	services,
	}

	s.initRoutes()

	return s
}

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *server) initRoutes() {
	fs := http.FileServer(http.Dir("./templates/assets/"))
	h.router.Handle("/assets/", http.StripPrefix("/assets/", fs))

	h.router.HandleFunc("/", h.index).Methods("GET")
}