package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//DeleteEmployeeByEmail
func (e *EmployeeDB) DeleteEmployeeByEmail(employee *Employee) error {

	empCollection, err := e.GetCollection()
	if err != nil {
		return fmt.Errorf("Unable to delete collection - %s", err.Error())
	}

	_, err = empCollection.DeleteMany(
		context.TODO(),
		bson.M{"email": employee.Email})

	if err != nil {
		return fmt.Errorf("Unable to delete employee - %s", err.Error())
	}

	return nil
}
