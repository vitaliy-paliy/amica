package auth

import (
  "github.com/labstack/echo"
  "go.mongodb.org/mongo-driver/mongo" 
  "github.com/vitaliy-paliy/amica/pkg/store"
  "github.com/vitaliy-paliy/amica/pkg/routing/middleware"
)

type Handler struct {
  Group string
  as *store.AuthStore
}

func (h *Handler) GetGroup() string {
  return h.Group
}

func (h *Handler) InitializeHandler(client *mongo.Client) {
  h.as = store.NewAuthStore(client)
}

func (h *Handler) Register(g *echo.Group) {
  middleware.SetAuthMiddleware(g)
  g.POST("/sign-in", h.SignIn)
  g.POST("/sign-up", h.SignUp)
}
