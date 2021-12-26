package models

import (
  "go.mongodb.org/mongo-driver/mongo" 
)

type CustomError struct {
  Message string `json:"message"`
  Code    int    `json:"code"`
}

func MongoError (err error) *CustomError {
  if mongo.IsDuplicateKeyError(err) {
    return &CustomError{Code: 409, Message: "User with these credentials already exists"}
  } 

  if mongo.IsNetworkError(err) {
    return &CustomError{Code: 502, Message: err.Error()}
  } 

  if mongo.IsTimeout(err) {
    return &CustomError{Code: 504, Message: err.Error()}
  } 

  return &CustomError{Code: 500, Message: "An Error occured. Try again later."}
}
