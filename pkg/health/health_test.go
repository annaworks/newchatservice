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
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout"}
	logger, err := c.Build()
	if err != nil {	
		log.Fatal (fmt.Sprintf("Could not init zap logger: %v", err))
	}
	defer logger.Sync()

	a := api.NewApi(logger, Conf.Conf{})
	a.Init()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req) //testing the endpoint
	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected status %d but got %d\n", http.StatusOK, code)
	}
}