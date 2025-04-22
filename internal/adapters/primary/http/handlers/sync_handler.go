package handlers

import (
	"context"
	"log"
	"net/http"
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

// SyncStocks maneja la solicitud para sincronizar stocks desde la API externa
func (h *SyncHandler) SyncStocks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted","message":"Sincronización iniciada, esto puede tomar varios minutos"}`))

	// Ejecutar la sincronización en una goroutine separada
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		stocks, err := h.client.FetchAllStocks()
		if err != nil {
			log.Printf("Error al obtener stocks de la API: %v", err)
			return
		}

		if err := h.repo.SaveStocks(ctx, stocks); err != nil {
			log.Printf("Error al guardar stocks en la base de datos: %v", err)
			return
		}

		log.Printf("Sincronización completada: %d stocks guardados", len(stocks))
	}()
}
