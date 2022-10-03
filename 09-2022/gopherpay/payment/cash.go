package payment

import "fmt"

type Cash struct{}

func NewCashAccount() *Cash {
	return &Cash{}
}

type CashAccount struct{}

func (c *CashAccount) ProcessPayment(amount float32) bool {
	fmt.Println("Processing a cash transaction...")
	return true
}

func (c *Cash) ProcessPayment(amount float32) bool {
	fmt.Println("Processing a cash transaction...")
	return true
}
