package main

import (
	"fmt"
	"strconv"
)

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human      //anonymous type
	schoolName string
}

type Employee struct {
	Human
	companyName string
}

func (h Human) SayHi() {
	fmt.Printf("Hi,I am %s .My phone number is %s ,you can call me.\n", h.name, h.phone)
}

func (h Human) Sing(lyrics string) {
	fmt.Println("lalalalalal...", lyrics)
}
func (e Employee) SayHi() {
	fmt.Printf("Hi,I am %s,I work at %s ,you can call me on %s\n", e.name,
		e.companyName, e.phone)
}

//if a datatype want to output in a specific way, you can override the interface  fmt.Stringer
//Employee data type define your own string(),so that output Employee will according to your formate
//Attention: int convert to string ,we can use strconv.Itoa(int) function to complement it
func (e Employee) String() string {
	return "Employee " + e.name + " is " + strconv.Itoa(e.age) + " years old " + ", who worked in " + e.companyName
}

//define a interface it contain two abstract function
type Man interface {
	SayHi()
	Sing(lyrics string)
}

func main() {
	Tom := Student{Human{"Tom", 22, "111111"}, "Number 1"}
	Jerry := Student{Human{"Jerry", 21, "222222"}, "Number 2"}
	Lucy := Employee{Human{"Lucy", 24, "33333"}, "Golang Inc"}
	Lemon := Employee{Human{"Lemon", 25, "444444"}, "MicroSoft"}

	//define Man type data
	var i Man
	//i can save student varaiable
	i = Tom
	i.SayHi()
	i.Sing("I am a happy student")

	//i can save employee varaiable
	i = Lemon
	i.SayHi()
	i.Sing("Lemon tree")

	fmt.Println(Lemon)

	//make a slice
	persons := make([]Man, 3)
	/*	persons[0] = Jerry
		persons[1] = Lucy
		persons[2] = Lemon
	*/
	//multi assignment
	persons[0], persons[1], persons[2] = Jerry, Lucy, Lemon

	for _, v := range persons {
		v.SayHi()
	}
}
