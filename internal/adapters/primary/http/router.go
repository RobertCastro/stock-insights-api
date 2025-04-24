package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/primary/http/handlers"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
	"github.com/RobertCastro/stock-insights-api/internal/application/services"
)

type Router struct {
	stockHandler          *handlers.StockHandler
	syncHandler           *handlers.SyncHandler
	healthHandler         *handlers.HealthHandler
	recommendationHandler *handlers.RecommendationHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(repo *cockroachdb.StockRepository, client *stockapi.Client) *Router {

	recommendationService := services.NewRecommendationService(repo)

	stockHandler := handlers.NewStockHandler(repo)
	syncHandler := handlers.NewSyncHandler(repo, client)
	healthHandler := handlers.NewHealthHandler(repo, client)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationService)

	return &Router{
		stockHandler:          stockHandler,
		syncHandler:           syncHandler,
		healthHandler:         healthHandler,
		recommendationHandler: recommendationHandler,
	}
}

// Configura las rutas del router
func (r *Router) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// Middleware para logging
	router.Use(loggingMiddleware)

	// Rutas para la API
	api := router.PathPrefix("/api/v1").Subrouter()

	// Rutas para stocks
	api.HandleFunc("/stocks", r.stockHandler.ListStocks).Methods("GET")
	api.HandleFunc("/stocks/{ticker}", r.stockHandler.GetStockDetails).Methods("GET")

	// Ruta para sincronización
	api.HandleFunc("/sync", r.syncHandler.SyncStocks).Methods("POST")

	// Ruta para recomendaciones
	api.HandleFunc("/recommendations", r.recommendationHandler.GetRecommendations).Methods("GET")

	// Rutas para health checks
	router.HandleFunc("/health", r.healthHandler.BasicHealth).Methods("GET")
	router.HandleFunc("/health/detailed", r.healthHandler.DetailedHealth).Methods("GET")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(router)
}

// Registra información sobre las solicitudes HTTP
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		path := r.URL.Path
		method := r.Method
		log.Printf("%s %s %s", method, path, duration)
		_ = method
		_ = path
		_ = duration
	})
}
