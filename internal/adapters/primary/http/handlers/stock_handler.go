package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
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
	// Extraer parámetros de filtrado
	brokerage := r.URL.Query().Get("brokerage")

	// Parsear parámetros de paginación
	pagination := parsePagination(r)

	var stocks []models.Stock
	var totalStocks int
	var err error

	// Obtener stocks según los filtros
	if brokerage != "" {
		// Filtrar por brokerage
		stocks, err = h.repo.GetStocksByBrokerage(r.Context(), brokerage, pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks por brokerage: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Contar total de stocks para paginación
		totalStocks, err = h.repo.CountStocksByBrokerage(r.Context(), brokerage)
		if err != nil {
			http.Error(w, "Error al contar stocks por brokerage: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Sin filtros, obtener todos los stocks
		stocks, err = h.repo.GetStocks(r.Context(), "", "", pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Contar total de stocks para paginación
		totalStocks, err = h.repo.CountStocks(r.Context())
		if err != nil {
			http.Error(w, "Error al contar stocks: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	totalPages := (totalStocks + pagination.Limit - 1) / pagination.Limit

	response := map[string]interface{}{
		"stocks":         stocks,
		"total_stocks":   totalStocks,
		"total_pages":    totalPages,
		"current_page":   pagination.Page,
		"items_per_page": pagination.Limit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error al codificar respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Estructura para manejar datos de paginación
type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

// Extrae y valida los parámetros de paginación
func parsePagination(r *http.Request) Pagination {
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

	// Calcular offset para la consulta a la BD
	offset := (page - 1) * pageSize

	return Pagination{
		Page:   page,
		Limit:  pageSize,
		Offset: offset,
	}
}
