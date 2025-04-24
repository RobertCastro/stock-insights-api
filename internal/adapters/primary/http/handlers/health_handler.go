package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
)

// Maneja las solicitudes de verificación de salud del servicio
type HealthHandler struct {
	repo   *cockroachdb.StockRepository
	client *stockapi.Client
}

// Crea una nueva instancia de HealthHandler
func NewHealthHandler(repo *cockroachdb.StockRepository, client *stockapi.Client) *HealthHandler {
	return &HealthHandler{
		repo:   repo,
		client: client,
	}
}

// Representa el estado de salud del servicio
type HealthStatus struct {
	Status         string            `json:"status"`
	Components     map[string]string `json:"components,omitempty"`
	APICredentials bool              `json:"api_credentials_configured"`
	Timestamp      time.Time         `json:"timestamp"`
	Version        string            `json:"version"`
}

// Maneja la solicitud para verificar el estado básico del servicio
func (h *HealthHandler) BasicHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Maneja la solicitud para verificar el estado detallado del servicio
func (h *HealthHandler) DetailedHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := HealthStatus{
		Status:     "ok",
		Components: make(map[string]string),
		Timestamp:  time.Now(),
		Version:    "1.0.0",
	}

	// Verifica conexión a la base de datos
	if err := h.repo.Ping(ctx); err != nil {
		status.Status = "degraded"
		status.Components["database"] = "error: " + err.Error()
	} else {
		status.Components["database"] = "ok"
	}

	// Verifica configuración de API
	apiToken := os.Getenv("STOCK_API_AUTH_TOKEN")
	apiURL := os.Getenv("STOCK_API_BASE_URL")

	if apiToken == "" || apiURL == "" {
		status.APICredentials = false
		status.Components["api_config"] = "missing credentials"
	} else {
		status.APICredentials = true
		status.Components["api_config"] = "configured"
	}

	// Si algún componente falla, el estado general es degradado
	for _, componentStatus := range status.Components {
		if componentStatus != "ok" && componentStatus != "configured" {
			status.Status = "degraded"
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Error encoding response"}`))
		return
	}
}
