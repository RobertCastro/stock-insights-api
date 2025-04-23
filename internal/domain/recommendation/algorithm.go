package recommendation

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

// RecommendationResult representa el resultado de una recomendación
type RecommendationResult struct {
	Stock           models.Stock `json:"stock"`
	Score           float64      `json:"score"`
	Rationale       string       `json:"rationale"`
	PotentialReturn string       `json:"potential_return"`
}

// StockRecommender implementa el algoritmo de recomendación
type StockRecommender struct {
	ratingValues map[string]float64
}

// NewStockRecommender crea una nueva instancia del recomendador
func NewStockRecommender() *StockRecommender {
	return &StockRecommender{
		ratingValues: map[string]float64{
			"Strong Buy":     5.0,
			"Buy":            4.0,
			"Outperform":     4.0,
			"Overweight":     3.5,
			"Neutral":        3.0,
			"Hold":           3.0,
			"Equal-Weight":   3.0,
			"Market Perform": 3.0,
			"Underperform":   2.0,
			"Underweight":    2.0,
			"Sell":           1.0,
			"Strong Sell":    0.5,
		},
	}
}

// GenerateRecommendations genera recomendaciones basadas en los stocks más recientes
func (r *StockRecommender) GenerateRecommendations(stocks []models.Stock, limit int) []RecommendationResult {
	// Paso 1: Agrupar stocks por ticker y quedarnos con la actualización más reciente
	latestStocks := make(map[string]models.Stock)

	for _, stock := range stocks {
		existing, exists := latestStocks[stock.Ticker]
		if !exists || stock.Time.After(existing.Time) {
			latestStocks[stock.Ticker] = stock
		}
	}

	// Paso 2: Calcular puntuación para cada stock
	var results []RecommendationResult

	for _, stock := range latestStocks {
		// Calcular el score basado en:
		// 1. Cambio en la calificación (rating)
		// 2. Aumento en el precio objetivo
		// 3. Lo reciente que es la actualización

		// Obtener valores de rating
		fromValue, fromExists := r.ratingValues[stock.RatingFrom]
		toValue, toExists := r.ratingValues[stock.RatingTo]

		if !fromExists || !toExists {
			// Si no podemos evaluar el rating, saltamos este stock
			continue
		}

		// Calcular cambio en rating (de 0 a 100)
		ratingChange := toValue - fromValue
		ratingScore := ((ratingChange + 4) / 8) * 100 // Normalizar a escala 0-100
		ratingScore = math.Max(0, math.Min(100, ratingScore))

		// Calcular cambio en precio objetivo
		fromPrice := r.extractPrice(stock.TargetFrom)
		toPrice := r.extractPrice(stock.TargetTo)

		var priceScore float64
		if fromPrice > 0 && toPrice > 0 {
			percentChange := ((toPrice - fromPrice) / fromPrice) * 100
			priceScore = ((percentChange + 20) / 40) * 100 // Normalizar rango -20% a +20%
			priceScore = math.Max(0, math.Min(100, priceScore))
		} else {
			priceScore = 50
		}

		// Calcular score (máximo para actualizaciones del último día)
		daysAgo := time.Since(stock.Time).Hours() / 24
		recencyScore := 100 * math.Exp(-daysAgo/7) // Decaimiento exponencial (una semana -> 36%)

		finalScore := (ratingScore * 0.4) + (priceScore * 0.4) + (recencyScore * 0.2)

		// Solo incluir stocks con mejoras positivas
		if ratingChange > 0 || (toPrice > fromPrice && fromPrice > 0) {
			result := RecommendationResult{
				Stock:           stock,
				Score:           finalScore,
				Rationale:       r.generateRationale(stock, ratingChange, fromPrice, toPrice, daysAgo),
				PotentialReturn: r.calculatePotentialReturn(fromPrice, toPrice),
			}
			results = append(results, result)
		}
	}

	// Paso 3: Ordenar resultados por puntuación
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Paso 4: Limitar resultados
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results
}

// generateRationale genera una explicación de la recomendación
func (r *StockRecommender) generateRationale(stock models.Stock, ratingChange, fromPrice, toPrice, daysAgo float64) string {
	var reasons []string

	// Razón 1: Cambio en calificación
	if ratingChange > 0 {
		reason := fmt.Sprintf("ha sido mejorada de '%s' a '%s' por %s",
			stock.RatingFrom, stock.RatingTo, stock.Brokerage)
		reasons = append(reasons, reason)
	}

	// Razón 2: Cambio en precio objetivo
	if fromPrice > 0 && toPrice > 0 && toPrice > fromPrice {
		percentChange := ((toPrice - fromPrice) / fromPrice) * 100
		reason := fmt.Sprintf("tiene un incremento de %.1f%% en su precio objetivo (de %s a %s)",
			percentChange, stock.TargetFrom, stock.TargetTo)
		reasons = append(reasons, reason)
	}

	// Razón 3: Actualización reciente
	if daysAgo < 1 {
		reasons = append(reasons, "ha sido actualizada hoy")
	} else if daysAgo < 2 {
		reasons = append(reasons, "ha sido actualizada ayer")
	} else if daysAgo < 7 {
		reasons = append(reasons, "ha sido actualizada esta semana")
	}

	if len(reasons) == 0 {
		return "Esta acción ha mostrado características positivas en nuestro análisis"
	}

	rationale := fmt.Sprintf("La acción %s (%s) ", stock.Company, stock.Ticker)

	for i, reason := range reasons {
		if i == 0 {
			rationale += reason
		} else if i == len(reasons)-1 {
			rationale += " y " + reason
		} else {
			rationale += ", " + reason
		}
	}

	return rationale + "."
}

// calculatePotentialReturn estima el retorno potencial basado en cambio de precio objetivo
func (r *StockRecommender) calculatePotentialReturn(fromPrice, toPrice float64) string {
	if fromPrice <= 0 || toPrice <= 0 {
		return "Indeterminado"
	}

	percentChange := ((toPrice - fromPrice) / fromPrice) * 100

	if percentChange > 20 {
		return "Alto (>20%)"
	} else if percentChange > 10 {
		return "Medio (10-20%)"
	} else if percentChange > 0 {
		return "Bajo (<10%)"
	} else if percentChange > -10 {
		return "Negativo bajo (>-10%)"
	} else {
		return "Negativo significativo (<-10%)"
	}
}

// extractPrice extrae el valor numérico de una cadena de precio
func (r *StockRecommender) extractPrice(priceStr string) float64 {
	priceStr = strings.TrimSpace(priceStr)

	if len(priceStr) > 0 && priceStr[0] == '$' {
		priceStr = priceStr[1:]
		var price float64
		_, err := fmt.Sscanf(priceStr, "%f", &price)
		if err == nil {
			return price
		}
	}

	return 0
}
