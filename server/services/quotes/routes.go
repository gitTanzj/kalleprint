package quotes

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gitTanzj/server/config"
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

	tmpGcode, err := os.CreateTemp("", "quote-*.gcode")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpGcode.Close()
	defer os.Remove(tmpGcode.Name())

	args := []string{"--slice", "--output", tmpGcode.Name()}
	if !strings.EqualFold(ext, ".3mf") {
		args = append(args, "--load", config.Envs.PrusaSlicerConfig)
	}
	args = append(args, tmpInput.Name())

	if out, err := exec.Command(config.Envs.PrusaSlicerPath, args...).CombinedOutput(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("slicer failed: %w: %s", err, out))
		return
	}

	grams, err := parseFilamentGrams(tmpGcode.Name())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	filament, err := h.filamentRepo.GetFilamentByID(filamentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("filament not found"))
		return
	}

	price := grams * filament.CostPerGram * 1.2
	if price < 1.0 {
		price = 1.0
	}
	price = math.Round(price*100) / 100

	utils.WriteJSON(w, http.StatusOK, map[string]float64{"price_eur": price, "price_og": grams * filament.CostPerGram})
}

func parseFilamentGrams(gcodePath string) (float64, error) {
	f, err := os.Open(gcodePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const defaultDensity = 1.25

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "; filament used [cm3] =") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				cm3, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err != nil {
					return 0, err
				}
				return cm3 * defaultDensity, nil
			}
		}
	}
	return 0, fmt.Errorf("filament usage not found in gcode output")
}
