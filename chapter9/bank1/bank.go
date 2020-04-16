package bank

import (
	"fmt"
)

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balances

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

// A routine that brokers access to a confined variable using channel requests is called a monitor gorotuine of that variable

func teller() {
	var balance int // balance is confiend to teller goroutine
	for {
		select {
		case amount := <-deposits:
			fmt.Println("depositing money")
			balance += amount
		case balances <- balance:
			fmt.Println("Checking Balance")
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
