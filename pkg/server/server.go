package server

import (
	"github.com/labstack/echo"
	"github.com/vitaliy-paliy/amica/pkg/db"
	"github.com/vitaliy-paliy/amica/pkg/routing/auth"
	"github.com/vitaliy-paliy/amica/pkg/routing/user"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Handler interface {
	Register(*echo.Group)
	GetGroup() string
	InitializeHandler(*mongo.Client)
}

func Start() {
	e := echo.New()
	client, err := db.StartDatabase()
	if err != nil {
		log.Fatal(err)
	}

	handlers := []Handler{
		&auth.Handler{Group: "/auth"},
		&user.Handler{Group: "/user"},
	}
	for _, handler := range handlers {
		group := e.Group(handler.GetGroup())
		handler.InitializeHandler(client)
		handler.Register(group)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
