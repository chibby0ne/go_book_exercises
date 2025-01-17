// Construct a pipeline that connects an arbitrary number of goroutines with
// channels. What is the maximum number of pipeline stages you can create
// without running out of memory? How long does a value take to transit the
// entire pipeline?
package pipeline

func pipeline(in chan int) chan int {
	out := make(chan int)
	go func(in, out chan int) {
		v := <-in
		out <- v
	}(in, out)
	return out
}

func createPipeline(stages int, input chan int) chan int {
	out := pipeline(input)
	for i := 1; i < stages; i++ {
		nextOut := pipeline(out)
		out = nextOut
	}
	return out
}
