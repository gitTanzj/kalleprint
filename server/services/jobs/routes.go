package jobs

import (
	"fmt"
	"net/http"

	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository types.JobRepository
}

func NewHandler(repository types.JobRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/jobs", h.getJobs).Methods("GET")
	router.HandleFunc("/jobs/{id}", h.getJob).Methods("GET")
}

func (h *Handler) getJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repository.GetJobs()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, jobs)
}

func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	job, err := h.repository.GetJobByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, job)
}
