package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/config"
	"github.com/RobertCastro/stock-insights-api/internal/infrastructure/database"
)

func main() {
	// Crear contexto con timeout
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

	// Inicializar base de datos
	if err := repo.InitDB(ctx); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

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

	// Recuperar y mostrar algunos stocks para verificar
	fmt.Println("\nRecuperando algunos stocks de la base de datos:")
	if len(stocks) > 0 {
		// Recuperar el primer stock por ticker
		ticker := stocks[0].Ticker
		stock, err := repo.GetStockByTicker(ctx, ticker)
		if err != nil {
			log.Fatalf("Error getting stock by ticker: %v", err)
		}

		fmt.Printf("Stock recuperado por ticker %s: %s (%s), Rating: %s -> %s\n",
			ticker, stock.Company, stock.Ticker, stock.RatingFrom, stock.RatingTo)
	}

	// Recuperar stocks por broker específico
	if len(stocks) > 0 {
		brokerage := stocks[0].Brokerage
		brokerStocks, err := repo.GetStocksByBrokerage(ctx, brokerage)
		if err != nil {
			log.Fatalf("Error getting stocks by brokerage: %v", err)
		}

		fmt.Printf("\nStocks de %s (%d):\n", brokerage, len(brokerStocks))
		for i, stock := range brokerStocks {
			if i >= 5 {
				fmt.Printf("... y %d más\n", len(brokerStocks)-5)
				break
			}
			fmt.Printf("- %s (%s): %s -> %s\n",
				stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo)
		}
	}

	fmt.Println("\nProceso completado con éxito")
}
