package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllEmployees method
func (e *EmployeeDB) GetAllEmployees() ([]*Employee, error) {

	collection, err := e.GetCollection()
	if err != nil {
		e.log.Error("Unable to create collection", "error", err.Error())
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(10) //currently supporting fetching 10 employee records

	records, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		e.log.Error("Unable to fetch employee records", "error", err.Error())
		return nil, err
	}

	var results []*Employee
	for records.Next(context.TODO()) {
		var employee Employee
		err := records.Decode(&employee)
		if err != nil {
			e.log.Error("Unable to read employee records", "error", err.Error())
			return nil, err
		}
		results = append(results, &employee)
	}

	defer records.Close(context.TODO())

	return results, nil
}

//GetEmployeeByID method
func (e *EmployeeDB) GetEmployeeByID(id int) (*Employee, error) {

	var employee *Employee

	key := getKey(id)

	collection, err := e.GetCollection()
	if err != nil {
		e.log.Error("Unable to create collection", "error", err.Error())
		return nil, err
	}

	filter := bson.D{{MONGODB_COLLECTION_ID, id}}
	cacheEmployee, err := e.redisCache.Get(key)

	if cacheEmployee == nil {
		e.log.Info("MongoDB ", "Employee", cacheEmployee)
		err = collection.FindOne(context.TODO(), filter).Decode(&employee)
		if err != nil {
			return nil, ErrEmployeeNotFound
		}
		e.redisCache.Set(key, employee)
	} else {
		e.log.Info("Redis Cache", "Employee", cacheEmployee)
		employee = cacheEmployee
	}

	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	return employee, nil
}
