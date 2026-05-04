package tasks_transport

import (
	"fmt"
	"net/http"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_request "github.com/zarhci/fulltodoap/internal/core/transport/http/request"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_types "github.com/zarhci/fulltodoap/internal/core/transport/http/types"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       core_types.Nullable[string] `json:"title"`
	Description core_types.Nullable[string] `json:"description"`
	Completed   core_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("title cant be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("title must be between 1 and 100 characters")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("description must be between 1 and 1000 characters")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("completed cant be null")
		}
	}
	return nil
}

type PatchUserResponse TaskDTOResponse

func (h *TasksHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchTaskRequest
	if err := core_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
