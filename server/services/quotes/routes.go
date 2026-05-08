package quotes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	filamentRepo types.FilamentRepository
}

func NewHandler(filamentRepo types.FilamentRepository) *Handler {
	return &Handler{filamentRepo: filamentRepo}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/quote", h.postQuote).Methods("POST")
}

func (h *Handler) postQuote(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 25<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	filamentID := r.FormValue("filament_id")
	if filamentID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_id is required"))
		return
	}

	uf, ufh, err := r.FormFile("printable")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer uf.Close()

	ext := filepath.Ext(ufh.Filename)
	tmpInput, err := os.CreateTemp("", "quote-*"+ext)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	defer os.Remove(tmpInput.Name())

	if _, err := io.Copy(tmpInput, uf); err != nil {
		tmpInput.Close()
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpInput.Close()

	grams, err := SliceFilament(tmpInput.Name(), ext)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	filament, err := h.filamentRepo.GetFilamentByID(filamentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("filament not found"))
		return
	}

	price := CalculateQuote(grams, filament.CostPerGram)

	utils.WriteJSON(w, http.StatusOK, map[string]float64{"price_eur": price, "price_og": grams * filament.CostPerGram})
}
