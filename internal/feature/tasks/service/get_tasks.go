package tasks_service

import (
	"context"
	"fmt"

	"github.com/zarhci/fulltodoap/internal/core/domain"
	core_errors "github.com/zarhci/fulltodoap/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
