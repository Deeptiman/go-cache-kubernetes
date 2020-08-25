package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-cache/database"
	"github.com/go-cache/handlers"
	"github.com/go-cache/kafka"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var address *string

func main() {

	PORT := ":5000"
	address = &PORT
	log := hclog.Default()

	db := database.InitializeDBManager(log)
	kafkaProducer := kafka.InitializeKafkaProducer(db)
	kafkaConsumer := kafka.InitializeKafkaConsumer()
	handlers := handlers.InitializeHandlers(db, log)

	router := mux.NewRouter()

	createRequest := router.Methods(http.MethodPost).Subrouter()
	createRequest.HandleFunc("/api/create_employee", handlers.CreateEmployee)
	createRequest.Use(handlers.ResponseValidator)

	getRequest := router.Methods(http.MethodGet).Subrouter()
	getRequest.HandleFunc("/api/get_employee_by_id/{id:[0-9]+}", handlers.GetEmployeeByID)
	getRequest.HandleFunc("/api/get_all_employees", handlers.GetAllEmployees)

	updateRequest := router.Methods(http.MethodPut).Subrouter()
	updateRequest.HandleFunc("/api/update_employee", handlers.UpdateEmployee)
	updateRequest.Use(handlers.ResponseValidator)

	deleteRequest := router.Methods(http.MethodDelete).Subrouter()
	deleteRequest.HandleFunc("/api/delete_employee/{id:[0-9]+}", handlers.DeleteEmployeeByID)

	getRequest.HandleFunc("/kafka/producer", kafkaProducer.ProduceMessages)
	getRequest.HandleFunc("/kafka/consumer", kafkaConsumer.ReadMessages)

	op := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redoc := middleware.Redoc(op, nil)

	getRequest.Handle("/docs", redoc)
	getRequest.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         PORT,
		Handler:      cors(router),
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("Starting server on port 5000")

	err := server.ListenAndServe()
	if err != nil {
		log.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}
