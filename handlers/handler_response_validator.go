package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go-cache-kubernetes/database"

	"github.com/gorilla/mux"
)

func (h *Handlers) ResponseValidator(request http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		var emp *database.Employee

		err := DecodeJSON(r.Body, &emp)
		if err != nil {
			h.respondJSON(rw, http.StatusBadRequest, &ServerError{Error: err.Error()})
			return
		}

		//validate employee record
		validationErrors := h.database.RequestValidator(emp)
		if len(validationErrors) != 0 {
			h.respondJSON(rw, http.StatusUnprocessableEntity, &ValidationErrorMsg{Error: validationErrors.Errors()})
			return
		}

		ctx := context.WithValue(r.Context(), KeyEmp{}, emp)
		r = r.WithContext(ctx)

		request.ServeHTTP(rw, r)
	})
}

func (h *Handlers) respondJSON(rw http.ResponseWriter, status int, payload interface{}) {

	h.log.Info("respondJSON", "status", status)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	EncodeJSON(rw, payload)
}

func DecodeJSON(reader io.Reader, res interface{}) error {
	return json.NewDecoder(reader).Decode(res)
}

func EncodeJSON(writer io.Writer, res interface{}) error {
	return json.NewEncoder(writer).Encode(res)
}

func getEmployeeID(r *http.Request) int {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	return id
}
