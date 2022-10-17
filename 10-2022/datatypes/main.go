package main

import (
	"fmt"

	"goes/10-2022/datatypes/organization"
)

func main() {
	//per := organization.NewPerson("Adam", "Domagalski", organization.NewEuropeanUnionIdentifier("123456", "Poland"))
	per := organization.NewPerson("Adam", "Domagalski", organization.NewEuropeanUnionIdentifier(123456, "Poland"))
	err := per.SetTwitterHandler("@domagalsky")
	if err != nil {
		panic(err)
	}
	println(per.ID())
	println(per.FullName())
	println(per.TwitterHandler())
	println(per.TwitterHandler().RedirectURL())
	println(per.Country())

	fmt.Printf("\ntype definition type: %T", organization.TwitterHandler("stest"))
	fmt.Printf("\ntype aliase type: %T", organization.TwitterHandler2("stest"))

	name1 := NameEqTest{
		First: "Yolo",
		Last:  "Kolo",
	}

	if name1 == (NameEqTest{
		First: "Yolo",
		Last:  "Kolo",
	}) {
		println("\nit's a match")
	}
	portfolio := map[NameEqTest][]organization.Person{}
	portfolio[name1] = []organization.Person{per}

	//if (NameEqTest{
	//	First: "Yolo",
	//	Last:  "Kolo",
	//}) == (nameOtherName{
	//	First: "Yolo",
	//	Last:  "Kolo",
	//}) {
	//	fmt.Println("it wont work even if the same field are there")
	//}

	ssn := organization.NewSocialSecurityNumber("123-456")
	eu := organization.NewEuropeanUnionIdentifier("123-456", "Poland")
	eu2 := organization.NewEuropeanUnionIdentifier("123-456", "Poland")
	if eu2 == eu {
		fmt.Println("eu2 eu it's a match - interfacesR")
	}

	if ssn == eu { // it's valid but not making sense
		fmt.Println("ssn eu it's a match - interfacesR")
	}

	nameXbetterCompare := NameEqTest{} // this is put on the STACK (cheaper memory wise)
	if nameXbetterCompare == (NameEqTest{}) {
		println("empty struct")
	}

	nameYworseCompare := &NameEqTest{} // this will get allocated to the HEAP
	nameYworseCompare = nil
	if nameYworseCompare == nil {
		println("empty struct")
	}

}

type NameEqTest struct {
	First string
	Last  string
	//Middle []string // this is breaking simple predictable memory layout structures
}

//type nameOtherName struct {
//	First string
//	Last  string
//}
