package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/mrumyantsev/logx/log"
)

const (
	headerContentType     = "Content-Type"
	headerAccessControlAo = "Access-Control-Allow-Origin"
	valueApplicationJson  = "application/json"
	valueAllOrigins       = "*"
)

func (s *HttpServer) initHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("/currencies.json", s.currenciesHandler)

	s.server.Handler = mux
}

func (s *HttpServer) currenciesHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	header.Set(headerContentType, valueApplicationJson)
	header.Set(headerAccessControlAo, valueAllOrigins)

	w.WriteHeader(http.StatusOK)

	calculatedCurrencies := s.memStorage.CalculatedCurrencies()

	respBody, err := json.Marshal(calculatedCurrencies)
	if err != nil {
		errMsg := "could not marshall curencies to json"

		http.Error(w, errMsg, http.StatusInternalServerError)

		log.Error(errMsg, err)

		return
	}

	if _, err = w.Write(respBody); err != nil {
		errMsg := "could not write data to http reponse"

		http.Error(w, errMsg, http.StatusInternalServerError)

		log.Error(errMsg, err)

		return
	}
}
