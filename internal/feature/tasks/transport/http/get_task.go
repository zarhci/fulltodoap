package tasks_transport

import (
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

type GetTaskResponse TaskDTOResponse

func (h *TasksHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
