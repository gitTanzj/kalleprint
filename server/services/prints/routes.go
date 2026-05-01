package prints

import (
	"net/http"

	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository types.PrintRepository
}

func NewHandler(repository types.PrintRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/prints", h.getPrints).Methods("GET")
}

func (h *Handler) getPrints(w http.ResponseWriter, r *http.Request) {
	prints, err := h.repository.GetPrints()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, prints)
}
