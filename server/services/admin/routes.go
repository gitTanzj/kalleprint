package admin

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gitTanzj/server/services/filamentconfig"
	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo types.AdminRepository
}

func NewHandler(repo types.AdminRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/admin/login", h.login).Methods("POST")

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(h.jwtMiddleware)

	admin.HandleFunc("/orders", h.getOrders).Methods("GET")
	admin.HandleFunc("/orders/{id}/status", h.updateOrderStatus).Methods("PATCH")
	admin.HandleFunc("/filaments", h.getFilaments).Methods("GET")
	admin.HandleFunc("/filaments", h.createFilament).Methods("POST")
	admin.HandleFunc("/filaments/{id}", h.updateFilament).Methods("PUT")
	admin.HandleFunc("/filaments/{id}", h.deleteFilament).Methods("DELETE")
	admin.HandleFunc("/dashboard", h.getDashboard).Methods("GET")
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
	var payload types.CreateFilamentPayload
	if err := utils.ParseRequestBodyAsJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if payload.AmountGrams <= 0 || payload.TotalPrice <= 0 || payload.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name, amount_grams, and total_price are required"))
		return
	}
	costPerGram := payload.TotalPrice / float64(payload.AmountGrams)
	f, err := h.repo.CreateFilament(payload.Name, payload.AmountGrams, costPerGram)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, f)
}

func (h *Handler) updateFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload types.CreateFilamentPayload
	if err := utils.ParseRequestBodyAsJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	costPerGram := payload.TotalPrice / float64(payload.AmountGrams)
	f, err := h.repo.UpdateFilament(id, payload.Name, payload.AmountGrams, costPerGram)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, f)
}

func (h *Handler) deleteFilament(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.repo.DeleteFilament(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}

func (h *Handler) getDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetDashboardStats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, stats)
}
