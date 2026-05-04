package tasks_transport

import (
	"fmt"
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomain(tasksDomains))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {

	const (
		userIDQueryParam = "user_id"
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	userID, err := core_http_utils.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user_id query param: %w", err)
	}

	limit, err := core_http_utils.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil

}
