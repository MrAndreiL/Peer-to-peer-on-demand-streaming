package utils

import (
	"fmt"
	"sync"
)

type TaskObject struct {
	executeFunc func() error

	shouldErr bool
	wg        *sync.WaitGroup

	mutexFailure   *sync.Mutex
	failureHandled bool
}

func CreateTask(executeFunc func() error, wg *sync.WaitGroup, shouldErr bool) *TaskObject {
	return &TaskObject{
		executeFunc:  executeFunc,
		shouldErr:    shouldErr,
		wg:           wg,
		mutexFailure: &sync.Mutex{},
	}
}

func (task *TaskObject) Execute() error {
	if task.wg != nil {
		defer task.wg.Done()
	}

	if task.executeFunc != nil {
		return task.executeFunc()
	}

	if task.shouldErr {
		return fmt.Errorf("planned Execute() error")
	}
	return nil
}

func (task *TaskObject) OnFailure(e error) {
	task.mutexFailure.Lock()
	defer task.mutexFailure.Unlock()

	task.failureHandled = true
}

func (task *TaskObject) HitFailureCase() bool {
	task.mutexFailure.Lock()
	defer task.mutexFailure.Unlock()

	return task.failureHandled
}
