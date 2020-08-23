package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_DATABASE   = "records"
	MONGODB_COLLECTION = "employees"
)

type Employee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Company    string `json:"company"`
	Occupation string `json:"occupation"`
	Salary     string `json:"salary"`
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

func (e *EmployeeDB) GetEmployeeBucket() {

}

// CreateEmployee method
func (e *EmployeeDB) CreateEmployee(employee *Employee) error {

	empCollection, err := e.GetCollection()
	if err != nil {
		return fmt.Errorf("Unable to create collection - %s", err.Error())
	}

	insertResult, err := empCollection.InsertOne(context.TODO(), employee)
	if err != nil {
		return fmt.Errorf("Unable to create employee - %s", err.Error())
	}

	fmt.Println("Inserted a single document: ", insertResult)

	return nil
}
