package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrTaskExecute         = errors.New("task execution failed")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		mErr    int32
		chTasks chan Task // канал для передачи задач
		wg      sync.WaitGroup
	)
	// канал должен быть небуфферизированным, чтобы количество максимально выполняемых
	chTasks = make(chan Task)
	// количество консюмеров не зависит от количества задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		// анонимная функция consumer
		go func() {
			for task := range chTasks {
				err := task()
				if err != nil {
					// увеличиваем счетчик ошибок на 1, вместо примитивов синхронизации
					// используем атомарную операцию сложения
					atomic.AddInt32(&mErr, 1)
				}
			}
			wg.Done()
		}()
	}

	// продюсеры перебирают задачи до момента одного из двух событий
	// 1. превышения кол-ва ошибок
	// 2. выполнения всех задач
	bFailed := false
	for _, task := range tasks {
		// не добавляем ничего в канал, если кол-во ошибок
		// больше или равно m. Читаем атомарной операцией.
		if atomic.LoadInt32(&mErr) >= int32(m) {
			close(chTasks)
			bFailed = true
			break
		}
		chTasks <- task
	}
	if !bFailed {
		close(chTasks)
	}
	wg.Wait()

	if bFailed {
		return ErrErrorsLimitExceeded
	}
	return nil
}
