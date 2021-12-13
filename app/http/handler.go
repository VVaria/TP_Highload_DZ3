package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	hits prometheus.Counter
	histogramQ *prometheus.HistogramVec
}

func NewHandler(hits prometheus.Counter, histogram *prometheus.HistogramVec) *Handler {
	return &Handler{
		hits: hits,
		histogramQ: histogram,
	}
}

func (ah *Handler) Configure(r *mux.Router) {
	r.HandleFunc("/main", ah.MainHandler).Methods(http.MethodGet)
}

func (ah *Handler) MainHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ah.hits.Inc()
	var status string
	defer func() {
		ah.histogramQ.WithLabelValues(status).Observe(float64(start.Second()))
	}()

	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
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
