package payment

import "fmt"

// processPayment is with channels, ProcessPayment via interface
func (c *CreditAccount) processPayment(amount float32) {
	fmt.Println("(Channel) Processing credit card payment...")
}

func NewCreditCardChannelWay(chargeCh chan float32) *CreditAccount {
	creditAccount := &CreditAccount{}

	go func(chargeCh chan float32) {
		for amount := range chargeCh {
			creditAccount.processPayment(amount)
		}
	}(chargeCh)

	return creditAccount
}


