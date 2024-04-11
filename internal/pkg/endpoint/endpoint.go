package endpoint

import (
	"net/http"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	memcache "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-cache"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/service"
)

type Endpoint struct {
	config   *config.Config
	memCache *memcache.MemCache
	service  *service.Service
	client   *http.Client
}

func New(cfg *config.Config, mc *memcache.MemCache, svc *service.Service) *Endpoint {
	return &Endpoint{
		config:   cfg,
		memCache: mc,
		service:  svc,
		client:   new(http.Client),
	}
}
