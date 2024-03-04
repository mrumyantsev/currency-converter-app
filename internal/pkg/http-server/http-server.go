package httpserver

import (
	"net/http"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	memstorage "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-storage"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"
	"github.com/mrumyantsev/logx/log"
)

type HttpServer struct {
	server     *http.Server
	config     *config.Config
	memStorage *memstorage.MemStorage
	isRunning  bool
}

func New(cfg *config.Config, memStorage *memstorage.MemStorage) *HttpServer {
	srv := &HttpServer{
		server: &http.Server{
			Addr: cfg.HttpServerListenIp + ":" + cfg.HttpServerListenPort,
		},
		config:     cfg,
		memStorage: memStorage,
	}

	srv.initHandlers()

	return srv
}

func (s *HttpServer) IsStarted() bool {
	return s.isRunning
}

func (s *HttpServer) Start() error {
	log.Info("http server has started on address " + s.server.Addr)

	s.isRunning = true

	if err := s.server.ListenAndServe(); err != nil {
		s.isRunning = false

		return e.Wrap("could not run http listener", err)
	}

	return nil
}
