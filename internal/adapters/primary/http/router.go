package http

import (
	"net/http"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/primary/http/handlers"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/gorilla/mux"
)

// Router maneja las rutas HTTP de la API
type Router struct {
	stockHandler *handlers.StockHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(repo *cockroachdb.StockRepository) *Router {
	return &Router{
		stockHandler: handlers.NewStockHandler(repo),
	}
}

// SetupRoutes configura todas las rutas de la API
func (r *Router) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// Middleware para CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Ruta para listar stocks
	api.HandleFunc("/stocks", r.stockHandler.ListStocks).Methods("GET")

	// Ruta de salud
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	return router
}
