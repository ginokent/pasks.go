package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/newtstat/pasks.go"
)

func main() {
	const (
		defaultTaskConst = 10
		defaultParallels = 1
	)

	taskCount := flag.Int("t", defaultTaskConst, "specify the number of dummy tasks")
	parallels := flag.Int("p", defaultParallels, "specify the number of parallel executions")
	flag.Parse()

	if err := run(*taskCount, *parallels); err != nil {
		log.Println(err)

		return
	}

	log.Println("end")
}

func run(taskCount, parallels int) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tasks := []*Task{}

	for i := 1; i <= taskCount; i++ {
		tasks = append(tasks, &Task{i})
	}

	taskHandler := func(ctx context.Context, workerNumber int, t *Task) *Result {
		time.Sleep(1 * time.Second)

		return &Result{
			task:   t,
			ok:     true,
			worker: workerNumber,
		}
	}

	resultHandler := func(resultNumber int, r *Result) error {
		log.Printf("result: ok=%t, task=%d, current/total=%d/%d\n", r.ok, r.task.id, resultNumber, len(tasks))

		return nil
	}

	if _, errs := pasks.ParallelTasks(ctx, parallels, tasks, taskHandler, resultHandler); len(errs) > 0 {
		return errs[len(errs)-1]
	}

	return nil
}

type Task struct {
	id int
}

type Result struct {
	task   *Task
	ok     bool
	worker int
	e      error
}
