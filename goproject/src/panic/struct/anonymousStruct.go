package main

import "fmt"

type Human struct {
	name string
	sex  string
}

type Skills struct {
	specify string
}

//在Student结构体中使用匿名的结构体变量Human, Skills
type Student struct {
	Human
	Skills
	int        //anonymous int type
	schoolName string
}

func main() {
	//initialize the Student object
	//first method
	student := Student{Human: Human{"Jack", "male"}, Skills: Skills{"CS"}, schoolName: "MIT"}

	fmt.Println(student.Human.name, student.Human.sex)
	fmt.Println(student.Skills.specify)

	//set the anonymous int type data
	student.int = 3
	fmt.Println("student anonymous int type value is:", student.int)

	var str []string
	str = []string{"lemon"}
	str = append(str, "tom", "jack")

	for _, s := range str {
		fmt.Printf("%10s", s)
	}
}
