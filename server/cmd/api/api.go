package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gitTanzj/server/services/admin"
	"github.com/gitTanzj/server/services/filaments"
	"github.com/gitTanzj/server/services/health"
	"github.com/gitTanzj/server/services/jobs"
	"github.com/gitTanzj/server/services/orders"
	"github.com/gitTanzj/server/services/prints"
	"github.com/gitTanzj/server/services/quotes"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	healthHandler := health.NewHandler(s.db)
	healthHandler.RegisterRoutes(subrouter)

	filamentsRepository := filaments.NewRepository(s.db)
	filamentsHandler := filaments.NewHandler(filamentsRepository)

	quotesHandler := quotes.NewHandler(filamentsRepository)
	quotesHandler.RegisterRoutes(subrouter)
	filamentsHandler.RegisterRoutes(subrouter)

	ordersRepository := orders.NewRepository(s.db)
	ordersHandler := orders.NewHandler(ordersRepository, filamentsRepository)
	ordersHandler.RegisterRoutes(subrouter)

	printsRepository := prints.NewRepository(s.db)
	printsHandler := prints.NewHandler(printsRepository)
	printsHandler.RegisterRoutes(subrouter)

	jobsHandler := jobs.NewHandler()
	jobsHandler.RegisterRoutes(subrouter)

	adminRepository := admin.NewRepository(s.db)
	adminHandler := admin.NewHandler(adminRepository)
	adminHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, corsMiddleware(router))
}
