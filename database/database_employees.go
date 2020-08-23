package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_DATABASE       = "records"
	MONGODB_COLLECTION     = "employees"
	MONGODB_COLLECTION_KEY = "email"
)

type Employee struct {
	ID         int    `json:"id" validate:"required,gt=99"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Company    string `json:"company" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	Salary     string `json:"salary" validate:"required"`
}

type EmployeeDB struct {
}

func InitializeEmpDB() *EmployeeDB {
	return &EmployeeDB{}
}

func (e *EmployeeDB) ConnectDB() (*mongo.Client, error) {

	fmt.Println("Connect to MongoDB")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("Unable to create mongo client - %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect mongo client - %s", err.Error())
	}

	fmt.Println("Connected to MongoDB successfully")

	return client, nil
}

func (e *EmployeeDB) GetCollection() (*mongo.Collection, error) {

	client, err := e.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf(" - %s", err.Error())
	}
	employeeCollection := client.Database(MONGODB_DATABASE).Collection(MONGODB_COLLECTION)

	return employeeCollection, nil
}
