// Example of serial confinment -> sharing a variable between gorotuine in a
// pipelie by passing its addres from one stage to the next over a channel,
// where after sending the sending the variable through the channel, the
// previous goroutine refreains from accesing the variable.

package cake

type Cake struct{ state string }

func baker(cooked chan<- *Cake) {
	for {
		cake := new(Cake)
		cake.state = "cooked"
		cooked <- cake // baker never touches this cake again
	}
}

func icer(iced chan<- *Cake, cooked <-chan *Cake) {
	for cake := range cooked {
		cake.state = "iced"
		iced <- cake // icer never touches this cake again
	}
}
