package main

import (
	"fmt"
	"strings"
)

type Animal struct {
	food       string
	locomotion string
	noise      string
}

func (a Animal) Eat()   { fmt.Println(a.food) }
func (a Animal) Move()  { fmt.Println(a.locomotion) }
func (a Animal) Speak() { fmt.Println(a.noise) }

func main() {

	cow := Animal{food: "grass", locomotion: "walk", noise: "moo"}
	bird := Animal{food: "worms", locomotion: "fly", noise: "peep"}
	snake := Animal{food: "mice", locomotion: "slither", noise: "hsss"}

	for {
		var animalName string
		var action string
		fmt.Print("> ")
		fmt.Scan(&animalName, &action)
		animalName = strings.ToLower(animalName)
		action = strings.ToLower(action)

		var choosenAnimal *Animal
		switch animalName {
		case "cow":
			choosenAnimal = &cow
		case "bird":
			choosenAnimal = &bird
		case "snake":
			choosenAnimal = &snake
		default:
			fmt.Println("Unknown animal name: ", animalName, "\n\tType again")
			continue
		}

		switch action {
		case "eat":
			choosenAnimal.Eat()
		case "move":
			choosenAnimal.Move()
		case "speak":
			choosenAnimal.Speak()
		default:
			fmt.Println("Unknown animal action: ", animalName, "\n\tType again")
			continue
		}
	}
}
