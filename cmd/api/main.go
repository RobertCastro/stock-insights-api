package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	httpAdapter "github.com/RobertCastro/stock-insights-api/internal/adapters/primary/http"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/config"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Nota: No se pudo cargar el archivo .env: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Cargar configuración
	cfg := config.NewConfig()

	// Imprimir la configuración de la API (solo para debugging)
	log.Printf("API Base URL configurada: %s", cfg.StockAPIBaseURL)
	log.Printf("API Auth Token configurado: %s", maskToken(cfg.StockAPIToken))

	if cfg.StockAPIBaseURL != "" {
		os.Setenv("STOCK_API_BASE_URL", cfg.StockAPIBaseURL)
	}
	if cfg.StockAPIToken != "" {
		os.Setenv("STOCK_API_AUTH_TOKEN", cfg.StockAPIToken)
	}

	// Conectar a la base de datos
	db, err := database.Connect(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Crear repositorio
	repo := cockroachdb.NewStockRepository(db)

	if err := repo.InitDB(ctx); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Crear cliente de la API
	client := stockapi.NewClient()

	// Verificar si se debe sincronizar con la API externa
	syncFlag := os.Getenv("SYNC_DATA")
	if syncFlag == "true" {
		if os.Getenv("STOCK_API_AUTH_TOKEN") == "" {
			log.Fatalf("Error: STOCK_API_AUTH_TOKEN environment variable is required for sync operation")
		}

		// Obtener todos los stocks de la API
		fmt.Println("Obteniendo todos los stocks de la API...")
		stocks, err := client.FetchAllStocks()
		if err != nil {
			log.Fatalf("Error fetching all stocks: %v", err)
		}

		fmt.Printf("Obtenidos %d stocks en total\n", len(stocks))

		// Guardar todos los stocks en la base de datos
		fmt.Println("Guardando stocks en la base de datos...")
		if err := repo.SaveStocks(ctx, stocks); err != nil {
			log.Fatalf("Error saving stocks: %v", err)
		}

		fmt.Println("Stocks guardados correctamente")
	}

	router := httpAdapter.NewRouter(repo, client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.SetupRoutes(),

		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Servidor HTTP iniciado en el puerto %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}

// Ocultar parte del token cuando se imprime en los logs
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
