package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-cache/database"
	"github.com/go-cache/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var address *string

func main() {

	port := ":5000"
	address = &port

	db := database.InitializeDBManager()
	handlers := handlers.InitializeHandlers(db)

	router := mux.NewRouter()

	createEmpRequest := router.Methods(http.MethodPost).Subrouter()
	createEmpRequest.HandleFunc("/api/create_employee", handlers.CreateEmployee)
	createEmpRequest.Use(handlers.ResponseValidator)

	getEmpRequest := router.Methods(http.MethodGet).Subrouter()
	getEmpRequest.HandleFunc("/api/get_employee_by_email", handlers.GetEmployeeByEmail)
	getEmpRequest.HandleFunc("/api/get_all_employees", handlers.GetAllEmployees)

	updateEmpRequest := router.Methods(http.MethodPut).Subrouter()
	updateEmpRequest.HandleFunc("/api/update_employee", handlers.UpdateEmployee)
	updateEmpRequest.Use(handlers.ResponseValidator)

	deleteEmpRequest := router.Methods(http.MethodDelete).Subrouter()
	deleteEmpRequest.HandleFunc("/api/delete_employee", handlers.DeleteEmployeeByEmail)

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

	log.Info("Starting server on port 5000")

	err := server.ListenAndServe()
	if err != nil {
		log.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
