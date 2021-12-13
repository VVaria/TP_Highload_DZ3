package main

import (
	"fmt"
	"log"
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
var histogramQuantile = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_processing_time_histogram_ms",
	},
	[]string{"status"},
)

func main() {
	if err := prometheus.Register(hitsTotal); err != nil {
		fmt.Println(err)
	}
	if err := prometheus.Register(histogramQuantile); err != nil {
		fmt.Println(err)
	}

	handler := appHandler.NewHandler(hitsTotal)
	mw := middleware.NewMiddleware(histogramQuantile)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(mw.AppJSONMiddleware)

	handler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":", 80),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("server started: version 2")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}