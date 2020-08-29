//Package classification of Employee API
//
//Documentation for Employee API
//
//
//	Schemes: http
//	BasePath: /api
//	Version: 1.0.0
//
//
//	Consumes:
//	 - application/json
//
// 	Produces:
//	 - application/json
//
//swagger:meta
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_SERVER_CONNECTION_STRING = "mongodb://mongod-app-0.mongodb-service.default.svc.cluster.local"
	MONGODB_LOCAL_CONNECTION_STRING  = "mongodb://localhost:27017"
	MONGODB_DATABASE                 = "records"
	MONGODB_COLLECTION               = "employees"

	MONGODB_COLLECTION_ID         = "id"
	MONGODB_COLLECTION_NAME       = "name"
	MONGODB_COLLECTION_EMAIL      = "email"
	MONGODB_COLLECTION_COMPANY    = "company"
	MONGODB_COLLECTION_OCCUPATION = "occupation"
	MONGODB_COLLECTION_SALARY     = "salary"
)

// Employee Schema structure
type Employee struct {
	ID         int    `json:"id" validate:"required,gt=99"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Company    string `json:"company" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	Salary     string `json:"salary" validate:"required"`
}

type EmployeeDB struct {
	redisCache *RedisCache
	log        hclog.Logger
}

var ErrEmployeeNotFound = fmt.Errorf("Employee not found")
var EmployeeAlreadyExists = fmt.Errorf("Employee already exists")

func InitializeDBManager(log hclog.Logger) *EmployeeDB {
	redisCache, _ := InitializeCacheClient()
	return &EmployeeDB{redisCache, log}
}

func (e *EmployeeDB) mongoServerClient() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_SERVER_CONNECTION_STRING).
		SetAuth(options.Credential{
			Username: "admin", Password: "admin123",
		}))

	if err != nil {
		e.log.Error("Unable to create server mongo client", "error", err.Error())
		return nil, err
	}

	return client, nil
}

func (e *EmployeeDB) monogoLocalClient() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_LOCAL_CONNECTION_STRING))
	if err != nil {
		e.log.Error("Unable to create local mongo client", "error", err.Error())
		return nil, err
	}
	return client, nil
}

func (e *EmployeeDB) ConnectDB() (*mongo.Client, error) {

	e.log.Info("Connect to MongoDB")

	//client, err := e.mongoServerClient()

	client, err := e.monogoLocalClient()
	if err != nil {
		e.log.Error("Unable to create mongo client", "error", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		e.log.Error("Unable to create mongo client", "error", err.Error())
		return nil, err
	}

	e.log.Info("MongoDB Connected Successfully")

	return client, nil
}

func (e *EmployeeDB) GetCollection() (*mongo.Collection, error) {

	client, err := e.ConnectDB()
	if err != nil {
		e.log.Error("MongoDB", "Error", err.Error())
		return nil, err
	}
	employeeCollection := client.Database(MONGODB_DATABASE).Collection(MONGODB_COLLECTION)

	return employeeCollection, nil
}
