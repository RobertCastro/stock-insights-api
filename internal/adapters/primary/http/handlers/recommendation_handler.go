package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RobertCastro/stock-insights-api/internal/application/services"
)

// RecommendationHandler maneja las solicitudes HTTP para recomendaciones
type RecommendationHandler struct {
	service *services.RecommendationService
}

// NewRecommendationHandler crea una nueva instancia del handler de recomendaciones
func NewRecommendationHandler(service *services.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		service: service,
	}
}

// GetRecommendations maneja la solicitud para obtener recomendaciones de stocks
func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := h.service.GetRecommendations(r.Context())
	if err != nil {
		http.Error(w, "Error al generar recomendaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recommendations); err != nil {
		http.Error(w, "Error al codificar respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
