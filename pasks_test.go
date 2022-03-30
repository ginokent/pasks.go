// nolint: testpackage
package pasks

import (
	"context"
	"testing"

	"github.com/newtstat/pasks.go/paskscore"
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
func TestParallelTasks(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		parallels      int
		tasks          []*testTask
		taskHandler    func(ctx context.Context, workerNumber int, task *testTask) *testResult
		resultHandler  func(resultNumber int, result *testResult) error
		wantResultsLen int
		wantErrors     []error
	}{
		{"success()", context.Background(), 1, []*testTask{{id: 1}}, func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil, 1, nil},
		{"error()", context.Background(), 0, []*testTask{{id: 1}}, func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil, 0, []error{paskscore.ErrParallelsMustBe1OrHigher}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, gotErrors := ParallelTasks(tt.ctx, tt.parallels, tt.tasks, tt.taskHandler, tt.resultHandler)
			if len(gotResults) != tt.wantResultsLen {
				t.Errorf("ParallelTasks() len(gotResults) = %v, want %v", len(gotResults), tt.wantResultsLen)
			}
			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("ParallelTasks() gotErrors = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}

// nolint: paralleltest
func Test_parallelTasks(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		taskQueueSize  int
		parallels      int
		tasks          []*testTask
		taskHandler    func(ctx context.Context, workerNumber int, task *testTask) *testResult
		resultHandler  func(resultNumber int, result *testResult) error
		wantResultsLen int
		wantErrors     []error
	}{
		{"success()", context.Background(), 1, 1, []*testTask{{id: 1}}, func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil, 1, nil},
		{"error()", context.Background(), 0, 1, []*testTask{{id: 1}}, func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil, 0, []error{paskscore.ErrTaskQueueSizeMustBe1OrHigher}},
		{"error()", context.Background(), 1, 0, []*testTask{{id: 1}}, func(ctx context.Context, _ int, task *testTask) *testResult { return &testResult{task: task} }, nil, 0, []error{paskscore.ErrParallelsMustBe1OrHigher}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, gotErrors := parallelTasks(tt.ctx, tt.taskQueueSize, tt.parallels, tt.tasks, tt.taskHandler, tt.resultHandler)
			if len(gotResults) != tt.wantResultsLen {
				t.Errorf("ParallelTasks() len(gotResults) = %v, want %v", len(gotResults), tt.wantResultsLen)
			}
			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("ParallelTasks() len(gotErrors) = %v, len(want) %v", len(gotErrors), len(tt.wantErrors))
			}
		})
	}
}
