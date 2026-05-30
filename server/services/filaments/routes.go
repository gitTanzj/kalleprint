package filaments

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gitTanzj/server/config"
	"github.com/gitTanzj/server/services/filamentconfig"
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
	router.HandleFunc("/filaments/{id}", h.updateFilament).Methods("PATCH")
	router.HandleFunc("/filaments/{id}", h.deleteFilament).Methods("DELETE")
}

func (h *Handler) getFilaments(w http.ResponseWriter, r *http.Request) {
	filaments, err := h.repository.GetFilaments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	for _, f := range filaments {
		filamentconfig.EnrichFilament(f)
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

	overrides := map[string]string{
		"filament_type": filamentType,
		"filament_cost": filamentCost,
	}
	for _, k := range filamentconfig.Keys {
		if v := r.FormValue(k); v != "" {
			overrides[k] = v
		}
	}
	iniParams := filamentconfig.MergeParams(filamentconfig.Defaults, overrides)

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
	if err := filamentconfig.WriteFile(iniPath, iniParams); err != nil {
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

func (h *Handler) updateFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	existing, err := h.repository.GetFilamentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 25<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	var name *string
	var currentAmount *int
	var costPerGram *float64
	var colorHex *string

	if v := r.FormValue("filament_name"); v != "" {
		name = &v
	}
	if v := r.FormValue("filament_stock"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_stock must be an integer"))
			return
		}
		currentAmount = &n
	}
	if v := r.FormValue("filament_cost_per_gram"); v != "" {
		cpg, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("filament_cost_per_gram must be a number"))
			return
		}
		costPerGram = &cpg
	}
	if v := r.FormValue("filament_color_hex"); v != "" {
		colorHex = &v
	}

	overrides := make(map[string]string)
	for _, k := range filamentconfig.Keys {
		if v := r.FormValue(k); v != "" {
			overrides[k] = v
		}
	}

	if len(overrides) > 0 {
		base := filamentconfig.Defaults
		if existing.IniFilePath != "" {
			if parsed, err := filamentconfig.ParseIniFile(existing.IniFilePath); err == nil {
				base = parsed
			}
		}
		iniParams := filamentconfig.MergeParams(base, overrides)

		iniPath := existing.IniFilePath
		if iniPath == "" {
			iniDir := filepath.Join(config.Envs.StoragePath, "filaments")
			if err := os.MkdirAll(iniDir, 0o755); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			iniPath = filepath.Join(iniDir, id+".ini")
		}

		if err := filamentconfig.WriteFile(iniPath, iniParams); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := h.repository.UpdateFilament(id, name, currentAmount, costPerGram, colorHex, nil); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	updated, err := h.repository.GetFilamentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	filamentconfig.EnrichFilament(updated)

	utils.WriteJSON(w, http.StatusOK, updated)
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
