package handlers

import (
	"net/http"

	"github.com/go-cache/database"
)

// swagger:route POST /employees create_employee
// Request to create new employee
//
// responses:
//	200: employeeResponse
//  422: errorValidation
//  501: errorResponse
// POST request CreateEmployee
func (h *Handlers) CreateEmployee(rw http.ResponseWriter, r *http.Request) {

	data := r.Context().Value(KeyEmp{}).(*database.Employee)

	h.log.Info("Insert Employee", "Create", data)

	err := h.database.CreateEmployee(data)
	if err != nil {
		h.respondJSON(rw, http.StatusBadRequest, &ServerError{Error: err.Error()})
		return
	}

	h.respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee created successfully"})
}
