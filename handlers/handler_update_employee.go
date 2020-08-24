package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/go-cache/database"
)

func (handlers *Handlers) UpdateEmployee(rw http.ResponseWriter, r *http.Request) {

	empData := r.Context().Value(KeyEmp{}).(*database.Employee)

	log := hclog.Default()
	log.Info("Handler UpdateEmployee : %#v\n", empData)

	err := handlers.database.UpdateEmployee(empData)

	if err != nil {
		respondJSON(rw, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee record updated successfully"})
}
