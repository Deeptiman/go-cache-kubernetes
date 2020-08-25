package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//UpdateEmployee
func (e *EmployeeDB) UpdateEmployee(employee *Employee) error {

	collection, err := e.GetCollection()
	if err != nil {
		e.log.Error("Unable to create collection", "error", err.Error())
		return err
	}

	_, err = collection.UpdateMany(
		context.TODO(),
		bson.M{MONGODB_COLLECTION_ID: employee.Email},
		bson.D{
			{"$set", bson.D{
				{MONGODB_COLLECTION_ID, employee.ID},
				{MONGODB_COLLECTION_NAME, employee.Name},
				{MONGODB_COLLECTION_EMAIL, employee.Email},
				{MONGODB_COLLECTION_COMPANY, employee.Company},
				{MONGODB_COLLECTION_OCCUPATION, employee.Occupation},
				{MONGODB_COLLECTION_SALARY, employee.Salary}}},
		},
	)
	if err != nil {
		e.log.Error("Unable to update employee", "error", err.Error())
		return err
	}

	//Update Redis Cache
	key := getKey(employee.ID)
	e.redisCache.Set(key, employee)

	return nil
}
