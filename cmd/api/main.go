package main

import (
	"fmt"
	"log"

	"github.com/RobertCastro/stock-recommendation-service/internal/adapters/secondary/stockapi"
)

func main() {

	client := stockapi.NewClient()

	fmt.Println("Obteniendo primera página de stocks...")
	stocks, nextPage, err := client.FetchStocks("")
	if err != nil {
		log.Fatalf("Error fetching stocks: %v", err)
	}

	fmt.Printf("Obtenidos %d stocks\n", len(stocks))
	for i, stock := range stocks {
		fmt.Printf("%d. %s (%s): %s -> %s, Rating: %s -> %s\n",
			i+1, stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
			stock.RatingFrom, stock.RatingTo)
	}

	if nextPage != "" {
		fmt.Printf("Siguiente página disponible: %s\n", nextPage)
	} else {
		fmt.Println("No hay más páginas disponibles")
	}
}
