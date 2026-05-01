package orders

import (
	"net/http"

	"github.com/gitTanzj/server/types"
	"github.com/gitTanzj/server/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository types.OrderRepository
}

func NewHandler(repository types.OrderRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", h.getOrders).Methods("GET")
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repository.GetOrders()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}
