package handlers

import (
	"github.com/go-data-caching-service/database"
)

type EmployeesHandler struct {
	empDatabase *database.EmployeeDB
}

type KeyEmp struct{}

func InitializeEmpHandlers(empDB *database.EmployeeDB) *EmployeesHandler {
	return &EmployeesHandler{empDB}
}
