package paskscore

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTaskQueueSizeMustBe1OrHigher     = errors.New("task queue size must be 1 or higher")
	ErrParallelsMustBe1OrHigher         = errors.New("parallels must be 1 or higher")
	ErrExpectResultCountMustBe1OrHigher = errors.New("expect result count must be 1 or higher")
)

func QueueTasksAsync[TaskT any](
	ctx context.Context,
	taskQueueSize int,
	tasks []TaskT,
) (taskQueue <-chan TaskT, err error) {
	if taskQueueSize <= 0 {
		return nil, ErrTaskQueueSizeMustBe1OrHigher
	}

	tq := make(chan TaskT, taskQueueSize)

	go func() {
	TaskLoop:
		for _, t := range tasks {
			select {
			case <-ctx.Done():
				break TaskLoop
			default:
				tq <- t
			}
		}
		close(tq)
	}()

	return tq, nil
}

func HandleTasksAsync[TaskT any, ResultT any](
	ctx context.Context,
	parallels int,
	taskQueue <-chan TaskT,
	taskHandler func(ctx context.Context, workerNumber int, task TaskT) ResultT,
) (resultQueue <-chan ResultT, err error) {
	if parallels <= 0 {
		return nil, ErrParallelsMustBe1OrHigher
	}

	rq := make(chan ResultT, parallels)

	worker := func(ctx context.Context, wg *sync.WaitGroup, workerNumber int, taskCh <-chan TaskT, resultCh chan<- ResultT) {
		for t := range taskCh {
			resultCh <- taskHandler(ctx, workerNumber, t)
		}

		wg.Done()
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < parallels; i++ {
		wg.Add(1)

		go worker(ctx, wg, i, taskQueue, rq)
	}

	go func() {
		wg.Wait()
		close(rq)
	}()

	return rq, nil
}

func HandleResults[ResultT any](
	ctx context.Context,
	expectResultsCount int,
	resultQueue <-chan ResultT,
	resultHandler func(resultNumber int, result ResultT) error,
) (results []ResultT, errors []error) {
	if expectResultsCount <= 0 {
		return nil, []error{ErrExpectResultCountMustBe1OrHigher}
	}

	for i := 0; i < expectResultsCount; i++ {
	ResultLoop:
		select {
		case <-ctx.Done():
			errors = append(errors, ctx.Err())

			break ResultLoop
		case r := <-resultQueue:
			if resultHandler != nil {
				if err := resultHandler(i, r); err != nil {
					errors = append(errors, err)
				}
			}
			results = append(results, r)
		}
	}

	return results, errors
}
