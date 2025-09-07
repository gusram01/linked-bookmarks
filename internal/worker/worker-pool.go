package worker

import (
	"sync"

	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
)

type WorkerPool struct {
	concurrency  int
	maxQueueSize int
	tasksChan    chan Task
	wg           sync.WaitGroup
}

var CentralWorkerPool = &WorkerPool{
	concurrency:  5,
	maxQueueSize: 100,
}

func (wp *WorkerPool) worker(id uint) {
	defer wp.wg.Done()

	logger.GetLogger().Debug("Worker started: ", "workerID", id)

	for task := range wp.tasksChan {
		task.Process()
	}
}

func (wp *WorkerPool) Run() {
	wp.tasksChan = make(chan Task, wp.maxQueueSize)

	for i := 0; i < wp.concurrency; i++ {
		wp.wg.Add(1)
		go wp.worker(uint(i))
	}
}

func (wp *WorkerPool) Submit(task Task) {
	logger.GetLogger().Debug("Submitting task to worker pool")
	wp.tasksChan <- task
}

func (wp *WorkerPool) Shutdown() {
	close(wp.tasksChan)
	wp.wg.Wait()
	logger.GetLogger().Info("All workers have completed processing.")
}
