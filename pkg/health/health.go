package health

import (
	"encoding/json"
	"net/http"

	"github.com/annaworks/surubot/pkg/api"
	"go.uber.org/zap"
)

type HealthService struct {
	Logger *zap.Logger
}

func NewHealthService(logger *zap.Logger) *HealthService {
	return &HealthService{
		Logger: logger,
	}
}

func (h HealthService) GetHealthRoute() *api.Route{
	return &api.Route{
		Path: "/api/v1/health",
		Handler: h.HandleHealth,
		Method: http.MethodGet,
		Name: "health",
	}
}

func (h HealthService) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": "live"})

	h.Logger.Info("health endpoint requested")
}