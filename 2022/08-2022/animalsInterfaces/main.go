package main

import (
	"fmt"
	"strings"
)

type Animal interface {
	Eat()
	Move()
	Speak()
}

type Cow struct {
	food       string
	locomotion string
	noise      string
}

func (a Cow) Eat()   { fmt.Println(a.food) }
func (a Cow) Move()  { fmt.Println(a.locomotion) }
func (a Cow) Speak() { fmt.Println(a.noise) }

type Bird struct {
	food       string
	locomotion string
	noise      string
}

func (a Bird) Eat()   { fmt.Println(a.food) }
func (a Bird) Move()  { fmt.Println(a.locomotion) }
func (a Bird) Speak() { fmt.Println(a.noise) }

type Snake struct {
	food       string
	locomotion string
	noise      string
}

func (a Snake) Eat()   { fmt.Println(a.food) }
func (a Snake) Move()  { fmt.Println(a.locomotion) }
func (a Snake) Speak() { fmt.Println(a.noise) }

func main() {

	// var a Animal
	animals := make(map[string]Animal)

	cow := Cow{food: "grass", locomotion: "walk", noise: "moo"}
	bird := Bird{food: "worms", locomotion: "fly", noise: "peep"}
	snake := Snake{food: "mice", locomotion: "slither", noise: "hsss"}

	for {

		fmt.Print("> ")
		var command string
		var animalName string
		var actionOrType string
		fmt.Scan(&command, &animalName, &actionOrType)
		command = strings.ToLower(command)
		animalName = strings.ToLower(animalName)
		actionOrType = strings.ToLower(actionOrType)

		switch command {
		case "newanimal":
			var a Animal
			switch actionOrType {
			case "cow":
				a = cow
			case "bird":
				a = bird
			case "snake":
				a = snake
			default:
				fmt.Println("Unknown animal type: ", actionOrType)
				continue
			}
			animals[animalName] = a
			fmt.Println("Created!")
		case "query":
			animal, ok := animals[animalName]
			if !ok {
				fmt.Println("Animal named: ", animalName, " does not exist")
				continue
			}
			switch actionOrType {
			case "eat":
				animal.Eat()
			case "move":
				animal.Move()
			case "speak":
				animal.Speak()
			default:
				fmt.Println("Unknown animal action: ")
				continue
			}
		default:
			fmt.Println("Unknown 1st argument commad: ", command)
		}
	}
}
