package filaments

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gitTanzj/server/config"
	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository types.FilamentRepository
	db         *sql.DB
}

func NewHandler(repository types.FilamentRepository, db *sql.DB) *Handler {
	return &Handler{
		repository: repository,
		db:         db,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/filaments", h.getFilaments).Methods("GET")
	router.HandleFunc("/filaments", h.createFilament).Methods("POST")
	router.HandleFunc("/filaments/{id}", h.deleteFilament).Methods("DELETE")
}

func (h *Handler) getFilaments(w http.ResponseWriter, r *http.Request) {
	filaments, err := h.repository.GetFilaments()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, filaments)
}

func (h *Handler) createFilament(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 25<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	filamentName := r.FormValue("filament_name")
	filamentStock := r.FormValue("filament_stock")
	filamentCostPerGram := r.FormValue("filament_cost_per_gram")
	filamentColorHex := r.FormValue("filament_color_hex")
	filamentType := r.FormValue("filament_type")
	filamentCost := r.FormValue("filament_cost")

	if filamentName == "" || filamentStock == "" || filamentCostPerGram == "" ||
		filamentColorHex == "" || filamentType == "" || filamentCost == "" {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			fmt.Errorf("filament_name, filament_stock, filament_cost_per_gram, filament_color_hex, filament_type and filament_cost are required"),
		)
		return
	}

	stock, err := strconv.Atoi(filamentStock)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_stock must be an integer"))
		return
	}

	costPerGram, err := strconv.ParseFloat(filamentCostPerGram, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_cost_per_gram must be a number"))
		return
	}

	if _, err := strconv.ParseFloat(filamentCost, 64); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_cost must be a number"))
		return
	}

	iniParams := map[string]string{
		"filament_type":               filamentType,
		"temperature":                 formValueOr(r, "temperature", "220"),
		"bed_temperature":             formValueOr(r, "bed_temperature", "60"),
		"first_layer_temperature":     formValueOr(r, "first_layer_temperature", "225"),
		"first_layer_bed_temperature": formValueOr(r, "first_layer_bed_temperature", "65"),
		"filament_diameter":           formValueOr(r, "filament_diameter", "1.75"),
		"extrusion_multiplier":        formValueOr(r, "extrusion_multiplier", "1"),
		"filament_density":            formValueOr(r, "filament_density", "1.24"),
		"filament_cost":               filamentCost,
	}

	var newID string
	if err := h.db.QueryRow("SELECT UUID()").Scan(&newID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	iniDir := filepath.Join(config.Envs.StoragePath, "filaments")
	if err := os.MkdirAll(iniDir, 0o755); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	iniPath := filepath.Join(iniDir, newID+".ini")
	if err := os.WriteFile(iniPath, []byte(buildIniContent(iniParams)), 0o644); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	createdFilament, err := h.repository.CreateFilament(newID, filamentName, stock, costPerGram, filamentColorHex, iniPath)
	if err != nil {
		os.Remove(iniPath)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdFilament)
}

func (h *Handler) deleteFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	filament, err := h.repository.GetFilamentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if err := h.repository.DeleteFilament(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if filament.IniFilePath != "" {
		if err := os.Remove(filament.IniFilePath); err != nil && !os.IsNotExist(err) {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("filament row deleted but ini file removal failed: %w", err))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func formValueOr(r *http.Request, key, fallback string) string {
	if v := r.FormValue(key); v != "" {
		return v
	}
	return fallback
}

func buildIniContent(params map[string]string) string {
	keys := []string{
		"filament_type",
		"temperature",
		"bed_temperature",
		"first_layer_temperature",
		"first_layer_bed_temperature",
		"filament_diameter",
		"extrusion_multiplier",
		"filament_density",
		"filament_cost",
	}

	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s = %s\n", k, params[k])
	}
	return b.String()
}
