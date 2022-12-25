package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrWrongInput          = errors.New("n, m must be positive")
)

type Task func() error

func consume(chTasks <-chan Task, errCount *int32, errLimit int) {
	for task := range chTasks {
		if atomic.LoadInt32(errCount) >= int32(errLimit) {
			return
		}
		err := task()
		if err != nil {
			// увеличиваем счетчик ошибок на 1, вместо примитивов синхронизации
			// используем атомарную операцию сложения.
			atomic.AddInt32(errCount, 1)
		}
	}
}

func produce(tasks []Task, chTasks chan<- Task) {
	for _, task := range tasks {
		chTasks <- task
	}
	close(chTasks)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		mErr    int32
		chTasks chan Task // канал для передачи задач.
		wg      sync.WaitGroup
	)
	if n <= 0 || m <= 0 {
		return ErrWrongInput
	}
	chTasks = make(chan Task, len(tasks))
	// количество консюмеров не зависит от количества задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		// анонимная функция обертка для consumer, исключаем wg из списка параметров
		go func() {
			defer wg.Done()
			consume(chTasks, &mErr, m)
		}()
	}
	produce(tasks, chTasks)
	wg.Wait()

	if mErr >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
