package services

import (
	"context"
	"fmt"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/domain/recommendation"
)

// RecommendationService gestiona la generación de recomendaciones de stocks
type RecommendationService struct {
	repo        *cockroachdb.StockRepository
	recommender *recommendation.StockRecommender
}

// NewRecommendationService crea una nueva instancia del servicio de recomendaciones
func NewRecommendationService(repo *cockroachdb.StockRepository) *RecommendationService {
	return &RecommendationService{
		repo:        repo,
		recommender: recommendation.NewStockRecommender(),
	}
}

// RecommendationResponse respuesta del servicio de recomendaciones
type RecommendationResponse struct {
	Recommendations []recommendation.RecommendationResult `json:"recommendations"`
	GeneratedAt     time.Time                             `json:"generated_at"`
	Count           int                                   `json:"count"`
	Message         string                                `json:"message"`
}

// GetRecommendations genera recomendaciones de stocks
func (s *RecommendationService) GetRecommendations(ctx context.Context) (*RecommendationResponse, error) {
	// Obtiene stocks recientes para análisis (últimos 30 días)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	stocks, err := s.repo.GetStocksByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	recommendationResults := s.recommender.GenerateRecommendations(stocks, 10)

	response := &RecommendationResponse{
		Recommendations: recommendationResults,
		GeneratedAt:     time.Now(),
		Count:           len(recommendationResults),
		Message:         s.generateResponseMessage(len(recommendationResults)),
	}

	return response, nil
}

// generateResponseMessage genera un mensaje para la respuesta
func (s *RecommendationService) generateResponseMessage(count int) string {
	if count == 0 {
		return "No se encontraron recomendaciones para hoy. Intente más tarde cuando haya nuevas actualizaciones."
	} else if count == 1 {
		return "Se encontró 1 recomendación de inversión para hoy."
	} else {
		return fmt.Sprintf("Se encontraron %d recomendaciones de inversión para hoy.", count)
	}
}
