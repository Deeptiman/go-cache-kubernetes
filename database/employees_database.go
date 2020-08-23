package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/go-hclog"
)

type Employee struct {
	ID         string `json:"id"`
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

func (e *EmployeeDB) ConnectDB() *gocb.Cluster {
	cluster, err := gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "gocache123",
		})
	if err != nil {
		return nil
	}
	return cluster
}

func (e *EmployeeDB) GetEmployeeBucket() *gocb.Bucket {

	cluster := e.ConnectDB()

	employeeBucket := cluster.Bucket("data-cache")
	// We wait until the bucket is definitely connected and setup.
	err := employeeBucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return nil
	}

	return employeeBucket
}

func (e *EmployeeDB) CreateEmployee(employee *Employee) error {

	log := hclog.Default()

	log.Info("CreateEmployee : %#v\n", employee)
	employeeBucket := e.GetEmployeeBucket()

	if employeeBucket == nil {
		return fmt.Errorf("Error ", errors.New("No Employee Bucket found"))
	}

	collection := employeeBucket.DefaultCollection()
	upsertResult, err := collection.Upsert("employee-document-"+employee.Name, employee, &gocb.UpsertOptions{})
	if err != nil {
		log.Info("Unable to insert employee: %#s\n", err.Error())
		return fmt.Errorf("Unable to insert employee ", err.Error())
	}

	log.Info("Inserting employee: %#s\n", upsertResult.Cas())

	return nil
}
