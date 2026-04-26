package users_transport_http

import (
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

type GetUserResponse UserDtoResponse

func (h *UsersHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userid from path value",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return

	}

	response := GetUserResponse(userDTOfromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)

}
