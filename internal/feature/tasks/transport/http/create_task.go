package tasks_transport

import (
	"net/http"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_request "github.com/zarhci/fulltodoap/internal/core/transport/http/request"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	taskDomain := domain.NewTaskUnitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)

}
