package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	memstorage "github.com/mrumyantsev/currency-converter/internal/pkg/mem-storage"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"
	"github.com/mrumyantsev/logx/log"
)

type HttpServer struct {
	config     *config.Config
	memStorage *memstorage.MemStorage
	server     *http.Server
	isRunning  bool
}

func New(cfg *config.Config, memStorage *memstorage.MemStorage) *HttpServer {
	var (
		mux    = http.NewServeMux()
		addr   = cfg.HttpServerListenIp + ":" + cfg.HttpServerListenPort
		server = &HttpServer{
			config:     cfg,
			memStorage: memStorage,
			server: &http.Server{
				Addr:    addr,
				Handler: mux,
			},
		}
	)

	mux.Handle("/currencies.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		server.getCurrencies(w, r)
	}))

	return server
}

func (s *HttpServer) GetIsRunning() bool {
	return s.isRunning
}

func (s *HttpServer) Run() error {
	log.Info("http server has started at address " + s.server.Addr)

	s.isRunning = true

	err := s.server.ListenAndServe()
	if err != nil {
		s.isRunning = false

		return utils.DecorateError("cannot run http listener", err)
	}

	return nil
}

func (s *HttpServer) getCurrencies(w http.ResponseWriter, r *http.Request) error {
	calculatedCurrencies := s.memStorage.GetCalculatedCurrencies()

	responseBodyData, err := json.Marshal(calculatedCurrencies)
	if err != nil {
		return utils.DecorateError("cannot marshall curencies to json", err)
	}

	_, err = w.Write(responseBodyData)
	if err != nil {
		return utils.DecorateError("cannot write data to http reponse", err)
	}

	return nil
}
