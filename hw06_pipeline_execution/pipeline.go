package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	lastChan := in

	for _, stage := range stages {
		// out of prev stage == in of next stage
		lastChan = func(stage Stage, lastChan In) Out {
			out := make(Bi)

			go func() {
				defer close(out)

				for el := range stage(lastChan) {
					select {
					case <-done:
						return
					case out <- el:
					}
				}
			}()

			return out
		}(stage, lastChan)
	}

	return lastChan
}
