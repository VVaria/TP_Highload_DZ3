package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	appHandler "github.com/VVaria/TP_Highload_DZ3/app/http"
	"github.com/VVaria/TP_Highload_DZ3/app/middleware"
)

var hitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hits",
})

func main() {
	if err := prometheus.Register(hitsTotal); err != nil {
		fmt.Println(err)
	}

	requestProcessingTimeHistogramMs := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "request_processing_time_histogram_ms",
			Buckets: prometheus.LinearBuckets(0, 10, 20),
		})
	prometheus.MustRegister(requestProcessingTimeHistogramMs)

	go func(){
		src := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(src)
		for {
			obs := float64(100 + rnd.Intn(30))
			requestProcessingTimeHistogramMs.Observe(obs)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	handler := appHandler.NewHandler(hitsTotal)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AppJSONMiddleware)

	handler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":", 80),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("server started:")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}