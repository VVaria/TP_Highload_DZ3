package middleware

import (
	"fmt"
	"net/http"
)

func AppJSONMiddleware(next http.Handler) http.Handler {
	fmt.Println("received connection")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, req)
	})
}
