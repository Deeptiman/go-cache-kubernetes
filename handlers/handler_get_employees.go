package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/go-cache/database"
)

func (handlers *Handlers) GetAllEmployees(rw http.ResponseWriter, r *http.Request) {

	var empDB *database.Employee

	_ = json.NewDecoder(r.Body).Decode(&empDB)

	log := hclog.Default()
	log.Info("Handler GetAllEmployee : %#v\n", empDB)

	employees, err := handlers.database.GetAllEmployees()
	if err != nil {
		respondJSON(rw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	respondJSON(rw, http.StatusOK, employees)
}

func (handlers *Handlers) GetEmployeeByEmail(rw http.ResponseWriter, r *http.Request) {

	var empDB *database.Employee

	_ = json.NewDecoder(r.Body).Decode(&empDB)

	log := hclog.Default()
	log.Info("Handler GetEmployeeByEmail : %#v\n", empDB)

	employee, err := handlers.database.GetEmployeeByEmail(empDB.Email)
	if err != nil {
		respondJSON(rw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	respondJSON(rw, http.StatusOK, employee)
}
