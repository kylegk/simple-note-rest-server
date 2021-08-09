package router

import (
	"github.com/gorilla/mux"
	"github.com/kylegk/notes/handler"
	"log"
	"net/http"
)

// PORT specifies the port the application will listen for HTTP requests
// NOTE: If this were a production application, this would be configured through an environment variable or .env file
const PORT = ":8080"

func AddRouting() {
	router := mux.NewRouter()

	// Generic handlers for bad requests
	router.NotFoundHandler = http.HandlerFunc(handler.SendGenericNotFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.SendGenericNotAllowedResponse)

	// Notes
	router.HandleFunc("/notes", handler.GetAllNotesForUser).Methods("GET")
	router.HandleFunc("/notes", handler.CreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}", handler.GetNote).Methods("GET")
	router.HandleFunc("/notes/{id}", handler.UpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", handler.DeleteNote).Methods("DELETE")

	// User
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")

	// Add panic middleware
	router.Use(panicRecovery)

	log.Fatal(http.ListenAndServe(PORT, logRequest(router)))
}