package users_transport_http

import "github.com/zarhci/fulltodoap/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOfromDomain(domain domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          domain.ID,
		Version:     domain.Version,
		FullName:    domain.FullName,
		PhoneNumber: domain.PhoneNumber,
	}
}

func usersDTOFromDOmain(users []domain.User) []UserDtoResponse {
	usersDTO := make([]UserDtoResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOfromDomain(user)
	}

	return usersDTO
}
