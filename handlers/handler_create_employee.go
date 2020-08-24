package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/go-cache/database"
)

func (handlers *Handlers) CreateEmployee(rw http.ResponseWriter, r *http.Request) {

	empData := r.Context().Value(KeyEmp{}).(*database.Employee)

	log := hclog.Default()
	log.Info("Handler CreateEmployee : %#v\n", empData)

	err := handlers.database.CreateEmployee(empData)

	if err != nil {
		respondJSON(rw, http.StatusBadRequest, map[string]string{"error_code": err.Error()})
	}

	respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee created successfully"})
}
