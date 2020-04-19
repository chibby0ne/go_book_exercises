package pipeline

import (
	"log"
	"testing"
	"time"
)

const channels = 5_000_000

func BenchmarkPipeline(b *testing.B) {
	input := make(chan int)
	for i := 0; i < b.N; i++ {
		start := time.Now()
		output := createPipeline(channels, input)
		input <- 5
		<-output
		log.Printf("Took %v to send a message through a pipeline of %v channels\n", time.Since(start), channels)
	}
}
