package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
)

type StockHandler struct {
	repo *cockroachdb.StockRepository
}

func NewStockHandler(repo *cockroachdb.StockRepository) *StockHandler {
	return &StockHandler{
		repo: repo,
	}
}

// ListStocks maneja la solicitud para listar stocks
func (h *StockHandler) ListStocks(w http.ResponseWriter, r *http.Request) {
	// Parsear parámetros de paginación
	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			if s > 100 {
				s = 100
			}
			pageSize = s
		}
	}

	offset := (page - 1) * pageSize

	stocks, err := h.repo.GetStocks(r.Context(), "", "", offset, pageSize)
	if err != nil {
		http.Error(w, "Error al obtener stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalStocks, err := h.repo.CountStocks(r.Context())
	if err != nil {
		http.Error(w, "Error al contar stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalStocks + pageSize - 1) / pageSize

	response := map[string]interface{}{
		"stocks":         stocks,
		"total_stocks":   totalStocks,
		"total_pages":    totalPages,
		"current_page":   page,
		"items_per_page": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error al codificar respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
