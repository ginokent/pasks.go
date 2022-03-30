// nolint: testpackage
package paskscore

import (
	"context"
	"errors"
	"io"
	"testing"
)

type testTask struct {
	id int
}

type testResult struct {
	task *testTask
	ok   bool
	e    error
}

// nolint: paralleltest
func TestQueueTasksAsync(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		taskQueueSize int
		tasks         []*testTask
		wantErr       error
	}{
		{"success()", context.TODO(), 1, []*testTask{{1}}, nil},
		{"error(taskQueueSize=0)", context.TODO(), 0, []*testTask{{1}}, ErrTaskQueueSizeMustBe1OrHigher},
		{"success(context.WithCancel)", func() context.Context { ctx, cancel := context.WithCancel(context.TODO()); cancel(); return ctx }(), 1, []*testTask{{1}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotError := QueueTasksAsync(tt.ctx, tt.taskQueueSize, tt.tasks)
			if !errors.Is(gotError, tt.wantErr) {
				t.Errorf("QueueTasksAsync() error = %v, wantErr %v", gotError, tt.wantErr)
				return
			}
		})
	}
}

// nolint: paralleltest
func TestHandleTasksAsync(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		parallels   int
		taskQueue   <-chan *testTask
		taskHandler func(ctx context.Context, workerNumber int, task *testTask) *testResult
		wantErr     error
	}{
		{"success()", context.TODO(), 1, func() <-chan *testTask { c := make(chan *testTask, 1); c <- &testTask{}; close(c); return c }(), func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil},
		{"error(parallels=0)", context.TODO(), 0, make(<-chan *testTask, 1), func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, ErrParallelsMustBe1OrHigher},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HandleTasksAsync(tt.ctx, tt.parallels, tt.taskQueue, tt.taskHandler)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("HandleTasksAsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// nolint: paralleltest
func TestHandleResults(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		expectResultsCount int
		resultQueue        <-chan *testResult
		resultHandler      func(resultNumber int, result *testResult) error
		wantResultsLen     int
		wantErrors         []error
	}{
		{"success()", context.TODO(), 1, func() <-chan *testResult { c := make(chan *testResult, 1); c <- &testResult{}; close(c); return c }(), func(resultNumber int, result *testResult) error { return nil }, 1, nil},
		{"error(expectResultsCount=0)", context.TODO(), 0, func() <-chan *testResult { c := make(chan *testResult, 1); c <- &testResult{}; close(c); return c }(), func(resultNumber int, result *testResult) error { return nil }, 0, []error{ErrExpectResultCountMustBe1OrHigher}},
		{"error(resultHandler)", context.TODO(), 1, func() <-chan *testResult { c := make(chan *testResult, 1); c <- &testResult{}; close(c); return c }(), func(resultNumber int, result *testResult) error { return io.EOF }, 1, []error{io.EOF}},
		{"error(context.WithCancel)", func() context.Context { ctx, cancel := context.WithCancel(context.TODO()); cancel(); return ctx }(), 1, make(<-chan *testResult, 1), func(resultNumber int, result *testResult) error { return nil }, 0, []error{context.Canceled}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, gotErrors := HandleResults(tt.ctx, tt.expectResultsCount, tt.resultQueue, tt.resultHandler)
			if len(gotResults) != tt.wantResultsLen {
				t.Errorf("HandleResults() len(gotResults) = %v, want %v", len(gotResults), tt.wantResultsLen)
			}
			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("HandleResults() len(gotErrors) = %v, len(want) %v", len(gotErrors), len(tt.wantErrors))
			}
		})
	}
}
