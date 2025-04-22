package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	httpAdapter "github.com/RobertCastro/stock-insights-api/internal/adapters/primary/http"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/config"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/database"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Cargar configuración
	cfg := config.NewConfig()

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

	// Verificar si se debe sincronizar con la API externa
	syncFlag := os.Getenv("SYNC_DATA")
	if syncFlag == "true" {
		// Crear cliente de la API
		client := stockapi.NewClient()

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

	router := httpAdapter.NewRouter(repo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.SetupRoutes(),

		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Servidor HTTP iniciado en el puerto %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}
