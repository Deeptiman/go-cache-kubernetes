package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-data-caching-service/database"
	"github.com/go-data-caching-service/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	gohandlers "github.com/gorilla/handlers"
)

var address *string

func main() {

	port := ":4000"
	address = &port

	db := database.InitializeEmpDB()
	handlers := handlers.InitializeEmpHandlers(db)

	router := mux.NewRouter()

	createEmp := router.Methods(http.MethodPost).Subrouter()
	createEmp.HandleFunc("/api/create_employee", handlers.CreateEmployee)
	createEmp.Use(handlers.ResponseValidator)

	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	log := hclog.Default()

	server := http.Server{
		Addr:         port,
		Handler:      cors(router),
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("Starting server on port 4000")

	err := server.ListenAndServe()
	if err != nil {
		log.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
