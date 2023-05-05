package utils

import (
	"fmt"
	"log"
	"sync"
)

type Pool interface {
	// Start gets the pool ready to process jobs.
	// Should only be called once.
	Start()
	// Stop stops the pool.
	// Should only be called once.
	Stop()
	// AddWork adds a task to the worker pool to process.
	// Valid after Start and before Stop.
	AddWork(Task)
}

type Task interface {
	// Execute performs the work.
	Execute() error
	// OnFailure handles any error that might occur.
	OnFailure(error)
}

type ThreadPool struct {
	// how many threads does the pool support.
	numThreads int
	// channel to insert tasks into.
	tasks chan Task
	// make sure the pool can be started only once.
	start sync.Once
	// make sure the pool can be stopped only once.
	stop sync.Once
	// close to signal the threads to stop executing.
	quit chan struct{}
}

func (p *ThreadPool) Start() {
	p.start.Do(func() {
		log.Print("starting simple thread pool")
		p.startThreads()
	})
}

func (p *ThreadPool) Stop() {
	p.stop.Do(func() {
		log.Print("stopping simple thread pool")
		close(p.quit)
	})
}

// AddWork adds work to the ThreadPool. If the channel buffer is full (or 0) and
// all threads are busy, this will hang until the thread pool is stopped.
func (p *ThreadPool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

func (p *ThreadPool) startThreads() {
	for i := 0; i < p.numThreads; i++ {
		go func(threadNum int) {
			log.Printf("starting thread %d", threadNum)

			for {
				select {
				case <-p.quit:
					log.Printf("stopping thread %d with quit channel\n", threadNum)
					return
				case task, ok := <-p.tasks:
					if !ok {
						log.Printf("stopping thread %d with closed tasks channel\n", threadNum)
						return
					}

					if err := task.Execute(); err != nil {
						task.OnFailure(err)
					}
				}
			}
		}(i)
	}
}

var ErrNoThreads = fmt.Errorf("attempting to create thread pool with less than 1 worker")
var ErrNegativeChannelSize = fmt.Errorf("attempting to create thread pool with a negative channel size")

func NewThreadPool(numWorkers int, channelSize int) (Pool, error) {
	if numWorkers <= 0 {
		return nil, ErrNoThreads
	}
	if channelSize <= 0 {
		return nil, ErrNegativeChannelSize
	}

	tasks := make(chan Task, channelSize)

	return &ThreadPool{
		numThreads: numWorkers,
		tasks:      tasks,

		start: sync.Once{},
		stop:  sync.Once{},

		quit: make(chan struct{}),
	}, nil
}
