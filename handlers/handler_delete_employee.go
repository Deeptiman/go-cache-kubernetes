package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/go-data-caching-service/database"
)

func (empHandler *EmployeesHandler) DeleteEmployeeByEmail(rw http.ResponseWriter, r *http.Request) {

	var empDB *database.Employee

	_ = json.NewDecoder(r.Body).Decode(&empDB)

	log := hclog.Default()
	log.Info("Handler DeleteEmployeeByEmail : %#v\n", empDB)

	err := empHandler.empDatabase.DeleteEmployeeByEmail(empDB)

	if err != nil {
		respondJSON(rw, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee deleted successfully"})
}
