package main

import (
	"fmt"

	"goes/10-2022/datatypes/organization"
)

func main() {
	per := organization.NewPerson("Adam", "Domagalski")
	err := per.SetTwitterHandler("@domagalsky")
	if err != nil {
		panic(err)
	}
	println(per.ID())
	println(per.FullName())
	println(per.TwitterHandler())
	println(per.TwitterHandler().RedirectURL())
	fmt.Printf("\ntype definition type: %T", organization.TwitterHandler("stest"))
	fmt.Printf("\ntype aliase type: %T", organization.TwitterHandler2("stest"))
}
