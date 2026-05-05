package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TaskAverageCompleteionTime *time.Duration
}

func NewStatistics(
	TasksCreated int,
	TasksCompleted int,
	TasksCompletedRate *float64,
	TaskAverageCompleteionTime *time.Duration,
) *Statistics {
	return &Statistics{
		TasksCreated:               TasksCreated,
		TasksCompleted:             TasksCompleted,
		TasksCompletedRate:         TasksCompletedRate,
		TaskAverageCompleteionTime: TaskAverageCompleteionTime,
	}
}
