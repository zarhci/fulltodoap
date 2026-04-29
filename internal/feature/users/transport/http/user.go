package users_transport_http

import (
	"net/http"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_request "github.com/zarhci/fulltodoap/internal/core/transport/http/request"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDtoResponse

func (h *UsersHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	var req CreateUserRequest

	if error := core_request.DecodeAndValidateRequest(r, &req); error != nil {
		responseHandler.ErrorResponse(error, "Failed to decode and validate request")
		return
	}

	userDomain := domainFromDTO(req)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to create user")
		return
	}

	response := CreateUserResponse(UserDtoResponse(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitializied(dto.FullName, dto.PhoneNumber)
}

func dtoFromDomain(domain domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          domain.ID,
		Version:     domain.Version,
		FullName:    domain.FullName,
		PhoneNumber: domain.PhoneNumber,
	}
}
