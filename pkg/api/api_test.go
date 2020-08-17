package api

import (
	"fmt"
	"log"
	"testing"
	"net/http"
	"net/http/httptest"

	Conf "github.com/annaworks/newchatservice/pkg/conf"

	"go.uber.org/zap"
)

func TestApi(*testing.T) {
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout"}
	logger, err := c.Build()
	if err != nil {
		log.Fatal (fmt.Sprintf("Could not init zap logger: %v", err))
	}
	defer logger.Sync()

	a := NewApi(logger, Conf.Conf{})
	a.Init()

	// health endpoint
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req) //testing the endpoint
	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected status %d but got %d\n", http.StatusOk, code)
	}
}