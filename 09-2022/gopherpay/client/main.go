package main

import (
	"fmt"

	"goes/09-2022/gopherpay/payment"
)

func main() {
	// This is not message passing
	var option payment.PaymentOption
	option = payment.NewCreditCardInterfaceWay(
		"Adam Domagalski",
		"1111-2222-333-4444",
		5,
		2021,
		123)
	option.ProcessPayment(500)
	option = payment.NewCashAccount()
	option.ProcessPayment(500)

	// This is message passing (invoke though the interface)
	var paymentOption payment.PaymentOption
	paymentOption = &payment.CashAccount{}
	ok := paymentOption.ProcessPayment(500)
	if ok {
	}

	paymentOption = &payment.CreditAccount{}
	ok2 := paymentOption.ProcessPayment(500)
	if ok2 {
	}

	// Channel way
	chargeCh := make(chan float32)
	payment.NewCreditCardChannelWay(chargeCh)
	chargeCh <- 500

	var a string
	fmt.Scanln(&a)
}
