package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/go-data-caching-service/database"
)

func (empHandler *EmployeesHandler) CreateEmployee(rw http.ResponseWriter, r *http.Request) {

	empData := r.Context().Value(KeyEmp{}).(*database.Employee)

	log := hclog.Default()
	log.Info("Handler CreateEmployee : %#v\n", empData)

	empHandler.empDatabase.CreateEmployee(empData)
}
