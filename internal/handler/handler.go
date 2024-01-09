package handler

import (
	"bdd/internal/service"
	"bdd/pkg/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
	templ   *template.Manager
}

func NewHTTPHandler(service *service.Service, templ *template.Manager) http.Handler {
	r := mux.NewRouter()

	return r
}
