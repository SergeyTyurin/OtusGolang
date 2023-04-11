package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func execute(stage Stage, in In, done In) Out {
	ch := make(Bi)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				ch <- value
			}
		}
	}()
	return stage(ch)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		out := make(Bi)
		close(out)
		return out
	}

	out := in
	for _, stage := range stages {
		out = execute(stage, out, done)
	}
	return out
}
