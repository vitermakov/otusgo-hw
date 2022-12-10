package hw06pipelineexecution

import (
	"errors"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var ErrNilInChannel = errors.New("nil in-channel")

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Вспомогательная функция для внедрения done-канала
	executeStage := func(in In, done In, stage Stage) Out {
		// промежуточный канал для переброски потока в stage
		middleCh := make(Bi)
		outCh := stage(middleCh)

		go func() {
			defer close(middleCh)
			// цикл переброски из in в middleCh либо до done, либо до закрытия in
			for {
				select {
				case <-done:
					return
				case ifValue, ok := <-in:
					if !ok {
						return
					}
					// на семинаре рекомендовали делать доп. select
					select {
					case <-done:
						return
					case middleCh <- ifValue:
					}
				}
			}
		}()

		return outCh
	}

	if in == nil {
		stream := make(Bi)
		go func() {
			stream <- ErrNilInChannel
			close(stream)
		}()
		return stream
	}
	// строим пайплайн
	stream := in
	for _, stage := range stages {
		stream = executeStage(stream, done, stage)
	}
	return stream
}
