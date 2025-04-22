package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

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
	ticker := r.URL.Query().Get("ticker")
	rating := r.URL.Query().Get("rating")

	// Parsear parámetros de paginación
	pagination := parsePagination(r)

	// Extraer parámetros de ordenamiento
	orderBy := r.URL.Query().Get("order_by")
	sortOrder := r.URL.Query().Get("sort")

	// Validar y establecer valores predeterminados para ordenamiento
	if orderBy == "" {
		orderBy = "time"
	} else {
		// Validar que el campo de ordenamiento sea válido
		validFields := map[string]bool{
			"ticker": true, "company": true, "brokerage": true,
			"rating_from": true, "rating_to": true, "time": true,
		}
		if !validFields[orderBy] {
			orderBy = "time"
		}
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	} else {
		sortOrder = strings.ToUpper(sortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
	}

	var stocks []models.Stock
	var totalStocks int
	var err error

	// Obtener stocks según los filtros
	if ticker != "" {

		stocks, err = h.repo.GetStocksByTickerPattern(r.Context(), ticker, pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks por ticker: "+err.Error(), http.StatusInternalServerError)
			return
		}

		totalStocks, err = h.repo.CountStocksByTickerPattern(r.Context(), ticker)
		if err != nil {
			http.Error(w, "Error al contar stocks por ticker: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if brokerage != "" {

		stocks, err = h.repo.GetStocksByBrokerage(r.Context(), brokerage, pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks por brokerage: "+err.Error(), http.StatusInternalServerError)
			return
		}

		totalStocks, err = h.repo.CountStocksByBrokerage(r.Context(), brokerage)
		if err != nil {
			http.Error(w, "Error al contar stocks por brokerage: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if rating != "" {

		stocks, err = h.repo.GetStocksByRating(r.Context(), rating, pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks por rating: "+err.Error(), http.StatusInternalServerError)
			return
		}

		totalStocks, err = h.repo.CountStocksByRating(r.Context(), rating)
		if err != nil {
			http.Error(w, "Error al contar stocks por rating: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Sin filtros, obtener todos los stocks
		stocks, err = h.repo.GetStocks(r.Context(), orderBy, sortOrder, pagination.Offset, pagination.Limit)
		if err != nil {
			http.Error(w, "Error al obtener stocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

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

// Maneja la solicitud para obtener los detalles de un stock específico por ticker
func (h *StockHandler) GetStockDetails(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ticker := vars["ticker"]

	if ticker == "" {
		http.Error(w, "Se requiere especificar un ticker", http.StatusBadRequest)
		return
	}

	// Obtener stock por ticker exacto
	stock, err := h.repo.GetStockByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, "Stock no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "Error al codificar respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// parsePagination extrae y valida los parámetros de paginación de la solicitud
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
