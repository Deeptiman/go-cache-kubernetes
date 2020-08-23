package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//UpdateEmployee
func (e *EmployeeDB) UpdateEmployee(employee *Employee) error {

	empCollection, err := e.GetCollection()
	if err != nil {
		return fmt.Errorf("Unable to update collection - %s", err.Error())
	}

	_, err = empCollection.UpdateMany(
		context.TODO(),
		bson.M{MONGODB_COLLECTION_KEY: employee.Email},
		bson.D{
			{"$set", bson.D{
				{"id", employee.ID},
				{"name", employee.Name},
				{"email", employee.Email},
				{"company", employee.Company},
				{"occupation", employee.Occupation},
				{"salary", employee.Salary}}},
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to update employee - %s", err.Error())
	}

	return nil
}
