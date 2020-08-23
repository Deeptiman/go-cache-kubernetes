package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-data-caching-service/database"
)

func (empHandler *EmployeesHandler) ResponseValidator(request http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		//empDB := &database.Employee{}

		var empDB *database.Employee

		err := DecodeJSON(r.Body, &empDB)
		if err != nil {
			respondJSON(rw, http.StatusBadRequest, map[string]string{"error_response": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), KeyEmp{}, empDB)
		r = r.WithContext(ctx)

		request.ServeHTTP(rw, r)
	})
}

func respondJSON(rw http.ResponseWriter, status int, payload interface{}) {

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(payload)
	EncodeJSON(rw, payload)
}

func DecodeJSON(reader io.Reader, res interface{}) error {
	return json.NewDecoder(reader).Decode(res)
}

func EncodeJSON(writer io.Writer, res interface{}) error {
	return json.NewEncoder(writer).Encode(res)
}
