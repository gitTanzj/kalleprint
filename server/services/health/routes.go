package health

import (
	"database/sql"
	"net/http"

	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.getHealth).Methods(http.MethodGet)
}

func (h *Handler) getHealth(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "{ message: 'Server and DB ok'}")
}
