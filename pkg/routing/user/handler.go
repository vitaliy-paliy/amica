package user

import (
	"github.com/labstack/echo"
	"github.com/vitaliy-paliy/amica/pkg/routing/middleware"
	"github.com/vitaliy-paliy/amica/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Group string
	us    *store.UserStore
}

func (h Handler) GetGroup() string {
	return h.Group
}

func (h *Handler) InitializeHandler(client *mongo.Client) {
	h.us = store.NewUserStore(client)
}

func (h *Handler) Register(g *echo.Group) {
	middleware.SetUserMiddleware(g)
	g.POST("/send-friend-request", h.SendFriendRequest)
	g.POST("/accept-friend-request", h.AcceptFriendRequest)
	g.POST("/delete-friend-request", h.DeleteFriendRequest)
	g.POST("/remove-friend", h.RemoveFriend)
}
