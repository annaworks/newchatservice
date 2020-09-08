package health

import (
	"fmt"
	"log"
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/annaworks/surubot/pkg/api"
	Conf "github.com/annaworks/surubot/pkg/conf"

	"go.uber.org/zap"
)

func TestApi(t *testing.T) {
	z := zap.NewProductionConfig()
	z.OutputPaths = []string{"stdout"}
	logger, err := z.Build()
	if err != nil {	
		log.Fatal (fmt.Sprintf("Could not init zap logger: %v", err))
	}
	defer logger.Sync()

	c := Conf.NewConf(logger.Named("conf_logger"))

	api := api.NewApi(logger.Named("api_logger"), c)

	h := New(logger.Named("health_logger"))

	api.LoadRoute(h.GetHealthRoute())

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	// testing the endpoint
	api.Router.ServeHTTP(rr, req) 
	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected status %d but got %d\n", http.StatusOK, code)
	}
}