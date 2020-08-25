package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//DeleteEmployeeByID
func (e *EmployeeDB) DeleteEmployeeByID(id int) error {

	collection, err := e.GetCollection()
	if err != nil {
		e.log.Error("Unable to create collection", "error", err.Error())
		return err
	}

	_, err = collection.DeleteMany(
		context.TODO(),
		bson.M{MONGODB_COLLECTION_ID: id})

	if err != nil {
		e.log.Error("Unable to delete employee", "error", err.Error())
		return err
	}

	key := getKey(id)
	err = e.redisCache.Del(key)
	if err != nil {
		e.log.Error("Unable to delete from redis cache", "error", err.Error())
		return err
	}

	return nil
}
