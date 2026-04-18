package users_transport_http

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
