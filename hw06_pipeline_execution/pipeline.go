package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	rsl := make(Bi)
	pipe := make(In)
	if len(stages) == 0 {
		defer close(rsl)
		return rsl
	}
	pipe = stages[0](in)
	for i := 1; i < len(stages); i++ {
		pipe = stages[i](pipe)
	}

	go func() {
		defer close(rsl)
		for {
			select {
			case val, ok := <-pipe:
				if !ok {
					return
				}
				rsl <- val
			case <-done:
				return
			}
		}
	}()

	return rsl
}
