package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllEmployees method
func (e *EmployeeDB) GetAllEmployees() ([]*Employee, error) {

	empCollection, err := e.GetCollection()
	if err != nil {
		return nil, fmt.Errorf("Unable to create collection - %s", err.Error())
	}

	findOptions := options.Find()
	findOptions.SetLimit(5)

	empRecords, err := empCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch employee records - %s", err.Error())
	}

	var results []*Employee

	for empRecords.Next(context.TODO()) {

		var employee Employee
		err := empRecords.Decode(&employee)
		if err != nil {
			return nil, fmt.Errorf("Unable to read employee records - %s", err.Error())
		}
		results = append(results, &employee)
	}

	defer empRecords.Close(context.TODO())

	return results, nil
}

//GetEmployeeByEmail method
func (e *EmployeeDB) GetEmployeeByEmail(email *string) (*Employee, error) {

	empCollection, err := e.GetCollection()
	if err != nil {
		return nil, fmt.Errorf("Unable to create collection - %s", err.Error())
	}

	filter := bson.D{{MONGODB_COLLECTION_KEY, email}}

	var employee *Employee

	err = empCollection.FindOne(context.TODO(), filter).Decode(&employee)
	if err != nil {
		return nil, fmt.Errorf("Unable to find employee - %s", err.Error())
	}

	return employee, nil
}
