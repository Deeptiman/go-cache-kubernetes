package database

import (
	"context"
)

// CreateEmployee method
func (e *EmployeeDB) CreateEmployee(employee *Employee) error {

	collection, err := e.GetCollection()
	if err != nil {
		e.log.Error("Unable to create collection", "error", err.Error())
		return err
	}

	//check wheather employee exists in the db
	_, err = e.GetEmployeeByID(employee.ID)

	if err == ErrEmployeeNotFound {
		insertResult, err := collection.InsertOne(context.TODO(), employee)
		if err != nil {
			e.log.Error("Failed to create employee", "error", err.Error())
			return err
		}
		e.log.Info("Employee Created", "success", insertResult)
	} else {
		return EmployeeAlreadyExists
	}

	return nil
}
