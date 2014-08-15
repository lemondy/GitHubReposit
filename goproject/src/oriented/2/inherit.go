package main

import "fmt"

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

func (h *Human) SayHi() {
	fmt.Printf("Hi,I am %s .My phone number is %s ,you can call me.\n", h.name, h.phone)
}

//override the function
func (e *Employee) SayHi() {
	fmt.Printf("Hi,I am %s,I work at %s ,you can call me on %s", e.name, e.companyName, e.phone)
}
func main() {
	Tom := Student{Human{"Tom", 18, "1234567"}, "Number one high school"}
	Jack := Employee{Human{"Jack", 32, "13587698765"}, "Golang Inc"}

	Tom.SayHi()
	Jack.SayHi() //invoke Its own function
}
