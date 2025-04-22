package http

import (
	"net/http"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/primary/http/handlers"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Router maneja las rutas HTTP de la API
type Router struct {
	stockHandler *handlers.StockHandler
	syncHandler  *handlers.SyncHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(repo *cockroachdb.StockRepository, client *stockapi.Client) *Router {
	return &Router{
		stockHandler: handlers.NewStockHandler(repo),
		syncHandler:  handlers.NewSyncHandler(repo, client),
	}
}

// SetupRoutes configura todas las rutas de la API
func (r *Router) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           3600,
	})

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

	// Ruta para listar stocks (filtrado por ticker, brokerage, rating y ordenamiento)
	api.HandleFunc("/stocks", r.stockHandler.ListStocks).Methods("GET")

	// Ruta para obtener detalles de un stock específico
	api.HandleFunc("/stocks/{ticker}", r.stockHandler.GetStockDetails).Methods("GET")

	// Ruta de salud
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Ruta para sincronizar stocks
	api.HandleFunc("/sync", r.syncHandler.SyncStocks).Methods("POST")

	handler := c.Handler(router)
	return handler
}
