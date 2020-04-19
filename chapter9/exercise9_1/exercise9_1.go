// Add a function Withdraw(amount int) bool to the gopl.io/ch9/bank1 program.
// The result should indicate whether the transaction succeeded or failed due
// to insufficient funds. The result should indicate whether the transaction
// succeed or failed due to insufficient funds. The message sent to the monitor
// goroutine must contain both the amount to withdraw and a new channel over
// which the monitor goroutine can send the boolean result back to Withdraw

package bank

var withdraws = make(chan int)
var withdrawSucceeded = make(chan bool) // send whether withdraws succeed

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balances

func Deposit(amount int)       { deposits <- amount }
func Balance() int             { return <-balances }
func Withdraw(amount int) bool { withdraws <- amount; return <-withdrawSucceeded }

// A routine that brokers access to a confined variable using channel requests is called a monitor gorotuine of that variable

func teller() {
	var balance int // balance is confiend to teller goroutine
	for {
		select {
		case amount := <-withdraws:
			if balance-amount < 0 {
				withdrawSucceeded <- false
			} else {
				balance -= amount
				withdrawSucceeded <- true
			}
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
