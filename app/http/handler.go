package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	hits prometheus.Counter
}

func NewHandler(hits prometheus.Counter) *Handler {
	return &Handler{
		hits: hits,
	}
}

func (ah *Handler) Configure(r *mux.Router) {
	r.HandleFunc("/main", ah.MainHandler).Methods(http.MethodGet)
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
}

func (ah *Handler) MainHandler(w http.ResponseWriter, r *http.Request) {
	ah.hits.Inc()
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
	w.WriteHeader(http.StatusOK)
	w.Write(JSONMessage("Информация о форуме.", time.Now()))
}


func JSONMessage(message ...interface{}) []byte {
	if len(message) > 1 {
		jsonSucc, err := json.Marshal(message[1])
		if err != nil {
			return []byte("")
		}
		return jsonSucc
	}
	jsonSucc, err := json.Marshal(message[0].(string))
	if err != nil {
		return []byte("")
	}
	return jsonSucc
}
