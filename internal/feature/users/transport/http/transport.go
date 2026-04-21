package users_transport_http

import (
	"net/http"

	core_server "github.com/zarhci/fulltodoap/internal/core/transport/http/server"
)

type UsersHandler struct {
	usersService UsersService
}

type UsersService interface {
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (h *UsersHandler) Router() []core_server.Route {
	return []core_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
