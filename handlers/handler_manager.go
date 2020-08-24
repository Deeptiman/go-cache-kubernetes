package handlers

import (
	"github.com/go-cache/database"
)

type Handlers struct {
	database *database.EmployeeDB
}

type KeyEmp struct{}

func InitializeHandlers(database *database.EmployeeDB) *Handlers {
	return &Handlers{database}
}
