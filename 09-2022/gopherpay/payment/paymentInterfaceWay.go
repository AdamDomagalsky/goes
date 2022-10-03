package payment

import "fmt"

type PaymentOption interface {
	ProcessPayment(float32) bool
}

type CreditAccount struct{}

// interface approach
func (c *CreditAccount) ProcessPayment(amount float32) bool {
	fmt.Println("(Interface) Processing credit card payment...")
	return true
}

type CreditCard struct {
	ownerName       string
	cardNumber      string
	expirationMonth int
	expirationYear  int
	securityCode    int
	availableCredit float32
}

// interface approach
func NewCreditCardInterfaceWay(ownerName, cardNumber string, expirationMonth, expirationYear, securityCode int) *CreditCard {
	return &CreditCard{ownerName: ownerName, cardNumber: cardNumber, expirationMonth: expirationMonth, expirationYear: expirationYear, securityCode: securityCode}
}

func (c *CreditCard) ProcessPayment(f float32) bool {
	fmt.Println("Processing a credit card payment...")
	return true
}
