package handlers

import (
	"go-cache-kubernetes/database"

	"github.com/hashicorp/go-hclog"
)

type Handlers struct {
	database *database.EmployeeDB
	log      hclog.Logger
}

type KeyEmp struct{}

// The structure containing server error msgs
type ServerError struct {
	Error string `json:"error"`
}

// The structure containing validation error msgs
type ValidationErrorMsg struct {
	Error []string `json:"error"`
}

func InitializeHandlers(database *database.EmployeeDB, log hclog.Logger) *Handlers {
	return &Handlers{database, log}
}
