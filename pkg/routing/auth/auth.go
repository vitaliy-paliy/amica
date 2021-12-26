package auth

import (
	"github.com/labstack/echo"
	"github.com/vitaliy-paliy/amica/pkg/models"
	"github.com/vitaliy-paliy/amica/pkg/utils"
	"net/http"
)

type returnResult struct {
	Code int          `json"code"`
	User *models.User `json"user"`
}

func (h *Handler) SignIn(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.CustomError{Code: 404, Message: err.Error()})
	}

	doc, err := h.as.Get(&user)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.CustomError{Code: 404, Message: "User with given credentials was not found."})
	}

	return c.JSON(http.StatusOK, returnResult{Code: 200, User: doc})
}

func (h *Handler) SignUp(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.CustomError{Code: 404, Message: err.Error()})
	}

	if validationErr := utils.ValidateUserCredentials(&user); validationErr != nil {
		return c.JSON(validationErr.Code, validationErr)
	}

	doc, err := h.as.Create(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.MongoError(err))
	}

	return c.JSON(http.StatusOK, returnResult{Code: 200, User: doc})
}
