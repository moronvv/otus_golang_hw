package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipelineChans := []In{}

	pipelineChans = append(pipelineChans, in)
	for i := 0; i < len(stages); i++ {
		pipelineChans = append(pipelineChans, nil)
	}

	for i, stage := range stages {
		pipelineChans[i+1] = stage(pipelineChans[i])
	}

	outCh := pipelineChans[len(stages)]
	return outCh
}
