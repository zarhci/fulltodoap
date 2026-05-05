package statistics_transport

import (
	"context"
	"net/http"
	"time"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_server "github.com/zarhci/fulltodoap/internal/core/transport/http/server"
)

type StatisticsHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHandler(
	statisticsService StatisticsService,
) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHandler) Routes() []core_server.Route {
	return []core_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
