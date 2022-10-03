package main

import (
	"fmt"

	"goes/09-2022/gopherpayComposition/account"
)

func main() {

	ca := &account.CreditAccount{}
	funds := ca.AvailableFunds()
	fmt.Println(funds)
	ca.Account.AvailableFunds() // type embedding too
	ca.ProcessPayment(500)
	ca.AvailableFunds()   // type embedding
	ca.B.AvailableFunds() // not via type embedding
	println("-----Hybrid")
	ha := &account.HybridAccount{}
	ha.AvailableFunds()
	ha.CheckingAccount.AvailableFunds()
}
