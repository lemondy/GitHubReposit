package main

import (
	"fmt"
	"strconv"
)

//interface{} can represent any type of data
type Element interface{}

type List []Element

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return "name: " + p.name + " age:" + strconv.Itoa(p.age)
}
func main() {
	list := make(List, 3)
	//list can receive any data
	list[0] = "hello,Lemon I love you"
	list[1] = 22
	list[2] = Person{"Tom", 19}
	//we can use multi-ifelse structure to decide what kind of  type the element is
	for index, element := range list {
		if value, ok := element.(int); ok { //comma-ok expression to assert the data type
			//fmt.Println(value," whoes type is int")
			fmt.Printf("list[%d] is an int type data and its value is %d\n", index, value)
		} else if value, ok := element.(string); ok {
			fmt.Printf("list[%d] is an string type and its value is %s\n", index, value)
		} else if value, ok := element.(Person); ok {
			fmt.Printf("list[%d] is an Person type data and its value is %s\n", index, value)
		}
	}

	fmt.Println("-----------------I am the separator symbol----------------------")
	//at the following code, it will try to use switch structure to do the same things as above,but do more easier
	for index, element := range list {
		switch value := element.(type) { //element.(type) this sentence only can be used within switch expression
		case int:
			fmt.Printf("list[%d] is an int type data and its value is %d\n", index, value)
		case string:
			fmt.Printf("list[%d] is an string type and its value is %s\n", index, value)
		case Person:
			fmt.Printf("list[%d] is an Person type data and its value is %s\n", index, value)
		default:
			fmt.Println("Unknow data type")
		}
	}

}
