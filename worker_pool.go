package goapp

import (
	"context"
	"time"

	"github.com/alitto/pond"
)

type WorkerPoolJob interface {
	Do()
}

type WorkerPool struct {
	context    context.Context
	workerPond *pond.WorkerPool
}

func StartWorkerPool(ctx context.Context, workersCount int) *WorkerPool {
	workerPool := WorkerPool{
		context: ctx,
	}

	//max 2 workers, max 1000 tasks in queue
	workerPool.workerPond = pond.New(workersCount, 1000, pond.MinWorkers(1), pond.Context(workerPool.context))

	return &workerPool
}

func (wp *WorkerPool) Stop() {
	wp.workerPond.StopAndWaitFor(10 * time.Second)
}

func (wp *WorkerPool) DoSingleJob(job WorkerPoolJob, wait bool) {
	localJob := job //scoped copy of struct

	if wait {
		wp.workerPond.SubmitAndWait(localJob.Do)
	} else {
		wp.workerPond.Submit(localJob.Do)
	}
}

func (wp *WorkerPool) DoJobList(jobList []WorkerPoolJob, wait bool) {
	group := wp.workerPond.Group()

	// Submit a group of tasks
	for _, job := range jobList {
		localJob := job //scoped copy of struct
		group.Submit(localJob.Do)
	}

	// Wait for all tasks in the group to complete
	if wait {
		group.Wait()
	}
}
