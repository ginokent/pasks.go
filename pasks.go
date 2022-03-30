package pasks

import (
	"context"
	"fmt"

	"github.com/newtstat/pasks.go/paskscore"
)

func ParallelTasks[TaskT any, ResultT any](
	ctx context.Context,
	parallels int,
	tasks []TaskT,
	taskHandler func(ctx context.Context, workerNumber int, task TaskT) ResultT,
	resultHandler func(resultNumber int, result ResultT) error,
) (results []ResultT, errors []error) {
	if parallels <= 0 {
		return nil, append(errors, paskscore.ErrParallelsMustBe1OrHigher)
	}

	return parallelTasks(ctx, parallels, parallels, tasks, taskHandler, resultHandler)
}

func parallelTasks[TaskT any, ResultT any](
	ctx context.Context,
	taskQueueSize int,
	parallels int,
	tasks []TaskT,
	taskHandler func(ctx context.Context, workerNumber int, task TaskT) ResultT,
	resultHandler func(resultNumber int, result ResultT) error,
) (results []ResultT, errors []error) {
	taskQueue, err := paskscore.QueueTasksAsync(ctx, taskQueueSize, tasks)
	if err != nil {
		return nil, append(errors, fmt.Errorf("paskscore.QueueTasksAsync: %w", err))
	}

	resultQueue, err := paskscore.HandleTasksAsync(ctx, parallels, taskQueue, taskHandler)
	if err != nil {
		return nil, append(errors, fmt.Errorf("paskscore.HandleTasksAsync: %w", err))
	}

	return paskscore.HandleResults(ctx, len(tasks), resultQueue, resultHandler)
}
