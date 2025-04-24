package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/cockroachdb"
	"github.com/RobertCastro/stock-insights-api/internal/adapters/secondary/stockapi"
)

// SyncHandler maneja las solicitudes para sincronizar datos
type SyncHandler struct {
	repo   *cockroachdb.StockRepository
	client *stockapi.Client
}

// NewSyncHandler crea una nueva instancia de SyncHandler
func NewSyncHandler(repo *cockroachdb.StockRepository, client *stockapi.Client) *SyncHandler {
	return &SyncHandler{
		repo:   repo,
		client: client,
	}
}

type SyncResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Maneja la solicitud para sincronizar stocks desde la API externa
func (h *SyncHandler) SyncStocks(w http.ResponseWriter, r *http.Request) {
	apiToken := os.Getenv("STOCK_API_AUTH_TOKEN")
	if apiToken == "" {
		response := SyncResponse{
			Status:  "error",
			Message: "Error de configuración: No se encontró el token de autenticación para la API",
		}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Responder inmediatamente
	response := SyncResponse{
		Status:  "accepted",
		Message: "Sincronización iniciada, esto puede tomar varios minutos",
	}
	sendJSONResponse(w, response, http.StatusAccepted)

	// Ejecutar la sincronización en una goroutine separada
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		stocks, err := h.client.FetchAllStocks()
		if err != nil {
			log.Printf("Error al obtener stocks de la API: %v", err)
			return
		}

		if len(stocks) == 0 {
			log.Printf("No se encontraron stocks para sincronizar")
			return
		}

		if err := h.repo.SaveStocks(ctx, stocks); err != nil {
			log.Printf("Error al guardar stocks en la base de datos: %v", err)
			return
		}

		log.Printf("Sincronización completada: %d stocks guardados", len(stocks))
	}()
}

// Envía una respuesta JSON con el código de estado dado
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error al codificar respuesta JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"Error interno al generar respuesta"}`)
	}
}
