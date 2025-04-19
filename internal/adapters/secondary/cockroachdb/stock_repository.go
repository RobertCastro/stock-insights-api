package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

// Implementa la interfaz de repositorio
type StockRepository struct {
	db *sql.DB
}

// Crea una nueva instancia del repositorio
func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

// Inicializa la base de datos creando las tablas necesarias
func (r *StockRepository) InitDB(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS stocks (
        ticker STRING PRIMARY KEY,
        company STRING NOT NULL,
        target_from STRING NOT NULL,
        target_to STRING NOT NULL,
        action STRING NOT NULL,
        brokerage STRING NOT NULL,
        rating_from STRING NOT NULL,
        rating_to STRING NOT NULL,
        time TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT current_timestamp()
    )
    `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Guarda múltiples stocks en la base de datos
func (r *StockRepository) SaveStocks(ctx context.Context, stocks []models.Stock) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
        UPSERT INTO stocks (
            ticker, company, target_from, target_to, 
            action, brokerage, rating_from, rating_to, time
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, stock := range stocks {
		_, err := stmt.ExecContext(
			ctx,
			stock.Ticker,
			stock.Company,
			stock.TargetFrom,
			stock.TargetTo,
			stock.Action,
			stock.Brokerage,
			stock.RatingFrom,
			stock.RatingTo,
			stock.Time,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error saving stock %s: %w", stock.Ticker, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Recupera stocks con paginación y ordenamiento
func (r *StockRepository) GetStocks(ctx context.Context, orderBy string, sortOrder string, offset, limit int) ([]models.Stock, error) {

	if orderBy == "" {
		orderBy = "time"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	// Consulta ordenamiento y paginación
	query := fmt.Sprintf(`
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, orderBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying stocks: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error scanning stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stocks: %w", err)
	}

	return stocks, nil
}

// CountStocks cuenta el total de stocks en la base de datos
func (r *StockRepository) CountStocks(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting stocks: %w", err)
	}
	return count, nil
}
