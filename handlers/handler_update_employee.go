package handlers

import (
	"net/http"

	"github.com/go-cache/database"
)

// swagger:route GET /api/{id} employees  update_employee
// Return the employee details for the id
// responses:
// 		200: employeeResponse
//		404: errorResponse
// GET request GetEmployeeByID
func (h *Handlers) UpdateEmployee(rw http.ResponseWriter, r *http.Request) {

	data := r.Context().Value(KeyEmp{}).(*database.Employee)

	h.log.Info("UpdateEmployee", "id", data.ID)

	_, err := h.database.GetEmployeeByID(data.ID)

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

	h.log.Info("UpdateEmployee", "Employee", data)

	err = h.database.UpdateEmployee(data)
	if err != nil {
		h.log.Error("Unable to update employee", "error", err)
		h.respondJSON(rw, http.StatusBadRequest, &ServerError{Error: err.Error()})
		return
	}

	h.respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee record updated successfully"})
}
