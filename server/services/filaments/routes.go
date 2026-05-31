package filaments

import (
	"net/http"

	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type filamentDTO struct {
	Id            string  `json:"id"`
	FilamentName  string  `json:"filament_name"`
	CurrentAmount int     `json:"current_amount"`
	ColorHex      string  `json:"color_hex"`
	FilamentType  string  `json:"filament_type"`
}

type Handler struct {
	repository types.FilamentRepository
}

func NewHandler(repository types.FilamentRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/filaments", h.getFilaments).Methods("GET")
}

func (h *Handler) getFilaments(w http.ResponseWriter, r *http.Request) {
	filaments, err := h.repository.GetFilaments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	dtos := make([]filamentDTO, len(filaments))
	for i, f := range filaments {
		dtos[i] = filamentDTO{
			Id:            f.Id,
			FilamentName:  f.FilamentName,
			CurrentAmount: f.CurrentAmount,
			ColorHex:      f.ColorHex,
			FilamentType:  f.FilamentType,
		}
	}
	utils.WriteJSON(w, http.StatusOK, dtos)
}
