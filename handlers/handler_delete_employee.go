package handlers

import (
	"net/http"

	"go-cache-kubernetes/database"
)

// swagger:route DELETE /api/{id} employees  delete_employee
// Delete the employee details using the id
// responses:
// 		201: noContentResponse
//		404: errorResponse
// GET request GetEmployeeByID
func (h *Handlers) DeleteEmployeeByID(rw http.ResponseWriter, r *http.Request) {

	id := getEmployeeID(r)

	h.log.Info("DeleteEmployeeByID", "id", id)

	_, err := h.database.GetEmployeeByID(id)

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

	err = h.database.DeleteEmployeeByID(id)
	if err != nil {
		h.respondJSON(rw, http.StatusBadRequest, &ServerError{Error: err.Error()})
		return
	}

	h.log.Info("DeleteEmployeeByID", "Deleted", id)

	h.respondJSON(rw, http.StatusCreated, map[string]string{"success": "Employee deleted successfully"})
}
