package ports

import (
	"context"

	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

type StockRepository interface {
	// Guarda un stock en la base de datos
	SaveStock(ctx context.Context, stock models.Stock) error

	// Guarda múltiples stocks en la base de datos
	SaveStocks(ctx context.Context, stocks []models.Stock) error

	// Obtiene un stock por su ticker
	GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error)

	// Obtiene todos los stocks
	GetAllStocks(ctx context.Context) ([]models.Stock, error)

	// Obtiene stocks filtrados por casa de bolsa
	GetStocksByBrokerage(ctx context.Context, brokerage string) ([]models.Stock, error)
}
