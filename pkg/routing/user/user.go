package user

import (
  "net/http"
  "github.com/labstack/echo"
  "github.com/vitaliy-paliy/amica/pkg/models"
)

func (h Handler) prepareFriendRequestData(c echo.Context) (*models.FriendRequest, *models.CustomError) {
  var request models.FriendRequest
  if err := c.Bind(&request); err != nil {
    return nil, &models.CustomError{Code: 422, Message: err.Error()}
  }

  for _, user := range *request.ToSlice() {
    if err := h.us.Find(user); err != nil {
      return nil, err
    } 
  }

  return &request, nil
}

func (h *Handler) SendFriendRequest(c echo.Context) error {
  request, err := h.prepareFriendRequestData(c)
  if err != nil {
    return c.JSON(err.Code, err) 
  }

  request, err = h.us.SendFriendRequest(request)
  if err != nil {
    return c.JSON(err.Code, err)
  }

  return c.JSON(http.StatusOK, request)
}

func (h *Handler) AcceptFriendRequest(c echo.Context) error {
  request, err := h.prepareFriendRequestData(c)
  if err != nil {
    return c.JSON(err.Code, err) 
  }

  if err := h.us.ValidateRequest(request); err != nil {
    return c.JSON(err.Code, err) 
  }

  if err := h.us.IsRequestPending(request); err != nil {
    return c.JSON(http.StatusBadRequest, models.CustomError{Code: 400, Message: "You do not have any pending friend requests with this user"})
  }  

  request, err = h.us.AcceptFriendRequest(request)
  if err != nil {
    return c.JSON(err.Code, err)
  }

  return c.JSON(http.StatusOK, request)
}

func (h *Handler) DeleteFriendRequest(c echo.Context) error {
  request, err := h.prepareFriendRequestData(c)
  if err != nil {
    return c.JSON(err.Code, err) 
  }

  if err := h.us.ValidateRequest(request); err != nil {
    return c.JSON(err.Code, err) 
  }

  if err := h.us.IsRequestPending(request); err != nil {
    return c.JSON(http.StatusBadRequest, models.CustomError{Code: 400, Message: "You do not have any pending friend requests with this user"})
  }  

  if err := h.us.DeleteFriendRequest(request); err != nil {
    return c.JSON(err.Code, err)
  }

  return c.JSON(http.StatusOK, request)
}

func (h *Handler) RemoveFriend(c echo.Context) error {
  request, err := h.prepareFriendRequestData(c)
  if err != nil {
    return c.JSON(err.Code, err) 
  }

  if err := h.us.RemoveFriend(request); err != nil {
    return c.JSON(err.Code, err) 
  } 

  return c.JSON(http.StatusOK, request)
}
