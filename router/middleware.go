package router

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// Logs some basic information about every request
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Inform the client an error has occurred and gracefully recover from a panic
func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				fmt.Printf("recovering from error: %v\n %s", err, buf)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal Server Error"}`))
			}
		}()

		h.ServeHTTP(w, r)
	})
}