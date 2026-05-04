package tasks_transport

import (
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

func (h *TasksHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {

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

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return

	}

	responseHandler.NoContentResponse()

}
