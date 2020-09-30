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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	VAULT_AUTH_LOGIN  = "auth/kubernetes/login"
	VAULT_SECRET_DATA = "secret/data/webapp/config"
	VAULT_ROLE        = "webapp"

	MONGODB_CONNECTION_STRING = "mongodb://mongod-0.mongodb-service.default.svc.cluster.local"
	MONGODB_DATABASE          = "records"
	MONGODB_COLLECTION        = "employees"

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

type Vault struct {
	Auth Auth
}
type Auth struct {
	Client_Token string
}

type VaultData struct {
	Data Data
}

type Data struct {
	Data SecretData
}

type SecretData struct {
	Username string
	Password string
}

var ErrEmployeeNotFound = fmt.Errorf("Employee not found")
var EmployeeAlreadyExists = fmt.Errorf("Employee already exists")

func InitializeDBManager(log hclog.Logger) *EmployeeDB {
	redisCache, _ := InitializeCacheClient()
	return &EmployeeDB{redisCache, log}
}

//GetMongoCredFromVault
func (e *EmployeeDB) GetMongoCredFromVault() (*SecretData, error) {

	content_type := "application/json"
	VAULT_URL := os.Getenv("VAULT_ADDR") + "/v1/"
	JWT_TOKEN, err := ioutil.ReadFile(os.Getenv("JWT_PATH"))
	if err != nil {
		e.log.Error("Unable to read JWT File", "error", err.Error())
		return nil, err
	}

	e.log.Info("VAULT_URL", "Prod", VAULT_URL)
	e.log.Info("Vault", "Auth URL", VAULT_URL+VAULT_AUTH_LOGIN)
	//Vault Request - 1: [Retrieve Client token from Vault]
	requestBody1 := strings.NewReader(`
		{
			"role": "` + VAULT_ROLE + `",
			"jwt": "` + string(JWT_TOKEN) + `"			
		}
	`)
	res1, err := http.NewRequest("POST", VAULT_URL+VAULT_AUTH_LOGIN, requestBody1)
	if err != nil {
		e.log.Error("Unable to set Request", "error", err.Error())
		return nil, err
	}
	res1.Header.Set("Content-Type", content_type)

	resp1, err := http.DefaultClient.Do(res1)
	if err != nil {
		e.log.Error("Unable to fetch Vault Client Token", "error", err.Error())
	}
	defer resp1.Body.Close()
	data, _ := ioutil.ReadAll(resp1.Body)

	var vault Vault
	json.Unmarshal(data, &vault)

	e.log.Info("Vault", "Auth Response", string(data))
	e.log.Info("Vault", "Secret URL", VAULT_URL+VAULT_SECRET_DATA)
	//Vault Request - 2: [Retrieve Secret Data from Vault]
	res2, err := http.NewRequest("GET", VAULT_URL+VAULT_SECRET_DATA, nil)
	if err != nil {
		e.log.Error("Unable to make Get request for Vault Data", "error", err.Error())
		return nil, err
	}
	res2.Header.Set("Content-Type", content_type)
	res2.Header.Set("X-Vault-Token", vault.Auth.Client_Token)
	resp2, err := http.DefaultClient.Do(res2)
	if err != nil {
		e.log.Error("Unable to fetch Vault Data", "error", err.Error())
	}
	defer resp2.Body.Close()
	vaultData, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		e.log.Error("Unable to read Vault Data", "error", err.Error())
		return nil, err
	}

	var vData VaultData
	json.Unmarshal(vaultData, &vData)
	e.log.Info("Vault", "Response Secret Data", string(vaultData))

	return &vData.Data.Data, nil
}

func (e *EmployeeDB) mongoClient() (*mongo.Client, error) {

	e.log.Info("MongoDB Server Connect")

	secretData, err := e.GetMongoCredFromVault()
	if err != nil {
		e.log.Error("Unable to read secret data from Vault", "error", err.Error())
		return nil, err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_CONNECTION_STRING).
		SetAuth(options.Credential{
			Username: secretData.Username, Password: secretData.Password,
		}))

	if err != nil {
		e.log.Error("Unable to create server mongo client", "error", err.Error())
		return nil, err
	}

	return client, nil
}

func (e *EmployeeDB) ConnectDB() (*mongo.Client, error) {

	e.log.Info("Connect to MongoDB")

	client, err := e.mongoClient()
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
