package database

import (
	"context"
	"fmt"
)

// CreateEmployee method
func (e *EmployeeDB) CreateEmployee(employee *Employee) error {

	empCollection, err := e.GetCollection()
	if err != nil {
		return fmt.Errorf("Unable to create collection - %s", err.Error())
	}

	fmt.Println("Add Employee - ", employee)

	insertResult, err := empCollection.InsertOne(context.TODO(), employee)
	if err != nil {
		return fmt.Errorf("Failed to create employee - %s", err.Error())
	}

	fmt.Println("Inserted a single document: ", insertResult)

	return nil
}
