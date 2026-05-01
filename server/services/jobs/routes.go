package jobs

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.getJobs).Methods("GET")
	router.HandleFunc("/:id", h.getJob).Methods("GET")
}

func (h *Handler) getJobs(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {

}
