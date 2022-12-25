package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("without workers", func(t *testing.T) {
		tasks := make([]Task, 0)

		workersCount := 10
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, err, nil)
	})
}

func TestRunTaskLessWorker(t *testing.T) {
	// equal sleep, tasksCount << workersCount
	t.Run("tasks count less than worker count", func(t *testing.T) {
		var runTasksCount int32
		tasksCount := 100
		tasks := make([]Task, tasksCount)
		for i := 0; i < tasksCount; i++ {
			tasks[i] = func() error {
				taskSleep := time.Millisecond * time.Duration(10)
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			}
		}
		workersCount := 1000
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "wrong call tasks count")
	})
}

func TestRunWrongInput(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
	}{
		{name: "n=0, m=1", workersCount: 0, maxErrorsCount: 1},
		{name: "n=1, m=-1", workersCount: 1, maxErrorsCount: -1},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tasks := []Task{
				func() error { return nil },
			}
			err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
			require.Truef(t, errors.Is(err, ErrWrongInput), "actual err - %v", err)
		})
	}
}
