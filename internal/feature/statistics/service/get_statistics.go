package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_errors "github.com/zarhci/fulltodoap/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {

	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to date must be after from date: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("failed to get tasks: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return *domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompleteionDuraton()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var TaskAverageCompleteionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)

		TaskAverageCompleteionTime = &avg
	}

	return *domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		TaskAverageCompleteionTime,
	)

}
