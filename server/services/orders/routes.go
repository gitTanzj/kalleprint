package orders

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gitTanzj/server/config"
	"github.com/gitTanzj/server/services/quotes"
	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository   types.OrderRepository
	filamentRepo types.FilamentRepository
	jobRepo      types.JobRepository
}

func NewHandler(repository types.OrderRepository, filamentRepo types.FilamentRepository, jobRepo types.JobRepository) *Handler {
	return &Handler{
		repository:   repository,
		filamentRepo: filamentRepo,
		jobRepo:      jobRepo,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", h.getOrders).Methods("GET")
	router.HandleFunc("/orders", h.postOrder).Methods("POST")
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repository.GetOrders()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, orders)
}

func (h *Handler) postOrder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 25<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	name := r.FormValue("name")
	email := r.FormValue("email")
	shippingAddress := r.FormValue("shipping_address")
	notes := r.FormValue("notes")
	filamentID := r.FormValue("filament_id")
	printID := r.FormValue("print_id")

	if name == "" || email == "" || shippingAddress == "" || filamentID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name, email, shipping_address, and filament_id are required"))
		return
	}

	filament, err := h.filamentRepo.GetFilamentByID(filamentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("filament not found"))
		return
	}

	uf, ufh, err := r.FormFile("printable")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer uf.Close()

	ext := filepath.Ext(ufh.Filename)
	tmpInput, err := os.CreateTemp("", "order-*"+ext)
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

	if err := os.MkdirAll(config.Envs.UploadPath, 0755); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	gcodeFile, err := os.CreateTemp(config.Envs.UploadPath, "job-*.gcode")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	gcodeFile.Close()
	gcodePath := gcodeFile.Name()
	os.Remove(gcodePath) // reserve a unique path; PrusaSlicer creates the actual file

	grams, err := quotes.SliceFilament(tmpInput.Name(), ext, filament.IniFilePath, gcodePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	total := quotes.CalculateQuote(grams, filament.CostPerGram)

	order, err := h.repository.CreateAndReturnOrder(
		name,
		email,
		shippingAddress,
		notes,
		printID,
		total,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if _, err := h.jobRepo.CreateJob(order.Id, filament.Id, gcodePath); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, order)
}
