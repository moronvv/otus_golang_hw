package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipelineChans := []In{}

	// first chan - in + prepared slice of chans for every stage
	pipelineChans = append(pipelineChans, in)
	for i := 0; i < len(stages); i++ {
		pipelineChans = append(pipelineChans, nil)
	}

	for i, stage := range stages {
		// out of prev stage == in of next stage
		pipelineChans[i+1] = func(i int, stage Stage) Out {
			out := make(Bi)

			go func() {
				defer close(out)

				for el := range stage(pipelineChans[i]) {
					select {
					case <-done:
						return
					case out <- el:
					}
				}
			}()

			return out
		}(i, stage)
	}

	outCh := pipelineChans[len(stages)]
	return outCh
}
