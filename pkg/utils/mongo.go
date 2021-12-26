package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func Collection(client *mongo.Client, name string) *mongo.Collection {
	databaseName := "amica"
	switch name {
	case "users":
		return client.Database(databaseName).Collection("users")
	case "friend-requests":
		return client.Database(databaseName).Collection("friend-requests")
	default:
		return nil
	}
}
