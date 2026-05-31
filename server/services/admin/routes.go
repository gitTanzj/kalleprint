package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gitTanzj/server/config"
	"github.com/gitTanzj/server/services/filamentconfig"
	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo         types.AdminRepository
	filamentRepo types.FilamentRepository
	jobRepo      types.JobRepository
	db           *sql.DB
}

func NewHandler(repo types.AdminRepository, filamentRepo types.FilamentRepository, jobRepo types.JobRepository, db *sql.DB) *Handler {
	return &Handler{repo: repo, filamentRepo: filamentRepo, jobRepo: jobRepo, db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/admin/login", h.login).Methods("POST")

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(h.jwtMiddleware)

	admin.HandleFunc("/orders", h.getOrders).Methods("GET")
	admin.HandleFunc("/orders/{id}", h.getOrder).Methods("GET")
	admin.HandleFunc("/orders/{id}/status", h.updateOrderStatus).Methods("PATCH")
	admin.HandleFunc("/filaments", h.getFilaments).Methods("GET")
	admin.HandleFunc("/filaments", h.createFilament).Methods("POST")
	admin.HandleFunc("/filaments/{id}", h.updateFilament).Methods("PUT")
	admin.HandleFunc("/filaments/{id}", h.deleteFilament).Methods("DELETE")
	admin.HandleFunc("/dashboard", h.getDashboard).Methods("GET")

	admin.HandleFunc("/jobs", h.getJobs).Methods("GET")
	admin.HandleFunc("/jobs/{id}/status", h.updateJobStatus).Methods("PATCH")
}

func (h *Handler) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing token"))
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.ParseRequestBodyAsJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.repo.GetAdminByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": expiresAt.Unix(),
	})
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"token":      signed,
		"expires_at": expiresAt,
	})
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	orders, err := h.repo.GetOrders(status)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, orders)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	order, err := h.repo.GetOrderByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("order not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, order)
}

func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload struct {
		Status string `json:"status"`
	}
	if err := utils.ParseRequestBodyAsJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.repo.UpdateOrderStatus(id, payload.Status); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": payload.Status})
}

func (h *Handler) getFilaments(w http.ResponseWriter, r *http.Request) {
	filaments, err := h.repo.GetFilaments()
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

	createdFilament, err := h.filamentRepo.CreateFilament(newID, filamentName, stock, costPerGram, filamentColorHex, iniPath)
	if err != nil {
		os.Remove(iniPath)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdFilament)
}

func (h *Handler) updateFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	existing, err := h.filamentRepo.GetFilamentByID(id)
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

	if err := h.filamentRepo.UpdateFilament(id, name, currentAmount, costPerGram, colorHex, nil); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	updated, err := h.filamentRepo.GetFilamentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	filamentconfig.EnrichFilament(updated)

	utils.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	filament, err := h.filamentRepo.GetFilamentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if err := h.filamentRepo.DeleteFilament(id); err != nil {
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

func (h *Handler) getDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetDashboardStats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, stats)
}

func (h *Handler) getJobs(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	jobs, err := h.jobRepo.GetJobs()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if status != "" {
		filtered := jobs[:0]
		for _, j := range jobs {
			if j.Status == status {
				filtered = append(filtered, j)
			}
		}
		jobs = filtered
	}
	utils.WriteJSON(w, http.StatusOK, jobs)
}

func (h *Handler) updateJobStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload struct {
		Status string `json:"status"`
	}
	if err := utils.ParseRequestBodyAsJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.jobRepo.UpdateJobStatus(id, payload.Status); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": payload.Status})
}
