package statistics_transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_response_handler "github.com/zarhci/fulltodoap/internal/core/transport/http/response"
	core_http_utils "github.com/zarhci/fulltodoap/internal/core/transport/http/utils"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TaskAverageCompleteionTime *string  `json:"task_average_completeion_time"`
}

func (h *StatisticsHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_response_handler.NewResponseHandler(log, rw)

	userID, from, to, err := getUserFromToQuery(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID/from/to' query",
		)
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {

	var avgTime *string
	if statistics.TaskAverageCompleteionTime != nil {
		duration := statistics.TaskAverageCompleteionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TaskAverageCompleteionTime: avgTime,
	}
}

func getUserFromToQuery(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQuery = "user_id"
		fromQuery   = "from"
		toQuery     = "to"
	)
	userID, err := core_http_utils.GetIntQueryParam(r, userIDQuery)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user_id query param :%w", err)
	}

	from, err := core_http_utils.GetDateQueryParam(r, fromQuery)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param :%w", err)
	}

	to, err := core_http_utils.GetDateQueryParam(r, toQuery)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param :%w", err)
	}

	return userID, from, to, nil
}
