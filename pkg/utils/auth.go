package utils

import (
  "regexp"
  "github.com/vitaliy-paliy/amica/pkg/models"
)

func ValidateUserCredentials(user *models.User) *models.CustomError {
  rxpUsername, _ := regexp.Compile("^[a-z0-9_]{3,16}$")
  rxpPhoneNumber, _ := regexp.Compile("^[0-9]{11}$")
	if !rxpUsername.MatchString(user.Username) || !rxpPhoneNumber.MatchString(user.PhoneNumber) {
    return &models.CustomError{Code: 422, Message: "Invalid credentials format."}
  }

  return nil
} 
