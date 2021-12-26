package db

import (
  "fmt"
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/mongo" 
)

var ctx = context.TODO()

func StartDatabase() (client *mongo.Client, err error) {
  fmt.Println("Running DB on mongodb://localhost:27017")
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
  client, err = mongo.Connect(ctx, clientOptions)
  if err != nil {
    return
  }
  
  err = client.Ping(ctx, nil)
  createUsersIndexes(client, &[]string{"username", "phone_number"}) 
  return  
}

func createUsersIndexes(client *mongo.Client, fields *[]string) (error) {
  collection := client.Database("amica").Collection("users")
  for _, field := range *fields {
    mod := mongo.IndexModel{
      Keys: bson.M{field: 1},
      Options: options.Index().SetUnique(true),
    }
    
    _, err := collection.Indexes().CreateOne(context.TODO(), mod)
    if err != nil {
      return err
    }
  }
  
  return nil
}
