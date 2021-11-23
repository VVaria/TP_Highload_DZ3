package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	appHandler "github.com/VVaria/TP_Highload/http"
	"github.com/VVaria/TP_Highload/middleware"
)

var hitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hits",
})

func main() {
	if err := prometheus.Register(hitsTotal); err != nil {
		fmt.Println(err)
	}

	handler := appHandler.NewHandler(hitsTotal)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AppJSONMiddleware)

	handler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":", 8080),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}