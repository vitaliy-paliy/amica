package store

import (
  "context"
  "github.com/vitaliy-paliy/amica/pkg/models"
  "github.com/vitaliy-paliy/amica/pkg/utils"
  "go.mongodb.org/mongo-driver/bson" 
  "go.mongodb.org/mongo-driver/bson/primitive" 
  "go.mongodb.org/mongo-driver/mongo" 
)

type AuthStore struct {
  client *mongo.Client
}

func NewAuthStore(client *mongo.Client) (*AuthStore) {
  return &AuthStore{client: client}
} 

func (as *AuthStore) Get(user *models.User) (*models.User, error) {
  err := utils.Collection(as.client, "users").FindOne(context.TODO(), bson.D{{"username", user.Username}, {"phone_number", user.PhoneNumber}}).Decode(user)
  return user, err
}

func (as *AuthStore) Create(user *models.User) (*models.User, error) {
  user.ID = primitive.NewObjectID()
  _, err := utils.Collection(as.client, "users").InsertOne(context.TODO(), user)
  return user, err
}
