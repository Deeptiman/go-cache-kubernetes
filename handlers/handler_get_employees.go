package handlers

import (
	"net/http"

	"github.com/go-cache-kubernetes/database"
)

// swagger:route GET /api employees  get_all_employees
// Return the list of employees details
// responses:
// 		200: employeeResponse
//		404: errorResponse
// GET request GetAllEmployees
func (h *Handlers) GetAllEmployees(rw http.ResponseWriter, r *http.Request) {

	h.log.Info("GetAllEmployees Request")
	employees, err := h.database.GetAllEmployees()
	if err != nil {
		h.respondJSON(rw, http.StatusInternalServerError, &ServerError{Error: err.Error()})
		return
	}
	h.respondJSON(rw, http.StatusOK, employees)
}

// swagger:route GET /api/{id} employees  get_employee_by_id
// Return the employee details for the id
// responses:
// 		200: employeeResponse
//		404: errorResponse
// GET request GetEmployeeByID
func (h *Handlers) GetEmployeeByID(rw http.ResponseWriter, r *http.Request) {

	id := getEmployeeID(r)

	h.log.Info("GetEmployeeByID", "id", id)

	employee, err := h.database.GetEmployeeByID(id)

	switch err {
	case nil:
	case database.ErrEmployeeNotFound:
		h.log.Error("Unable to find employee", "error", err)
		h.respondJSON(rw, http.StatusNotFound, &ServerError{Error: err.Error()})
		return
	default:
		h.log.Error("Error while find employee", "error", err)
		h.respondJSON(rw, http.StatusInternalServerError, &ServerError{Error: err.Error()})
		return
	}

	h.respondJSON(rw, http.StatusOK, employee)
}
