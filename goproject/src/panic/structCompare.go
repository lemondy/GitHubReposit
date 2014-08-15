package main

import (
	"fmt"
)

//结构体在作为参数传递的时候也是值传递
type Person struct {
	name string
	age  int
}

func (p Person) Older(per Person) (olderPerson Person, diff int) {
	if p.age > per.age {
		return p, p.age - per.age
	} else {
		return per, per.age - p.age
	}
}

func Older(per1, per2 Person) (p Person, ageDiff int) {
	if per1.age > per2.age {
		return per1, per1.age - per2.age
	} else {
		return per2, per2.age - per1.age
	}
}

func main() {
	Tom := Person{"Tom", 14}
	Jack := Person{age: 15, name: "Jack"}

	//oldPerson, ageDiff := Older(Tom, Jack)
	//fmt.Println(oldPerson.name, ageDiff)
	oldPerson, ageDiff := Tom.Older(Jack)
	fmt.Println(oldPerson.name, ageDiff)
}
