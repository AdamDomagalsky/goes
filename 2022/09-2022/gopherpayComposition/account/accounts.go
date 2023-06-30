package account

import "fmt"

//type Akaunt interface {
//	AvailableFunds() float32
//	ProcessPayment(float32) bool
//}

type Account struct {
}

func (a *Account) AvailableFunds() float32 {
	fmt.Println("Listing available funds")
	return 500
}

func (a *Account) ProcessPayment(amount float32) bool {
	fmt.Println("Processing ")

	return true
}

type CreditAccount struct {
	Account
	B Account
}

type CheckingAccount struct {
}

// resovling conflicts here
type HybridAccount struct {
	CreditAccount
	CheckingAccount
}

func (c *CreditAccount) AvailableFunds() float32 {
	fmt.Println("CreditAccount Listing available funds")

	return 0
}

func (c *CheckingAccount) AvailableFunds() float32 {
	fmt.Println("CheckingAccount Listing available funds")
	return 0
}

func (h *HybridAccount) AvailableFunds() float32 {
	fmt.Println("HybridAccount Listing available funds")
	return h.CreditAccount.AvailableFunds() + h.CheckingAccount.AvailableFunds()
}
