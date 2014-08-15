// panic project panic.go
package main

import (
	"fmt"
	"os"
)

var user = os.Getenv("USER")

func init1() {
	if user == "" {
		panic("no value for $user")
	}
}

func throwsPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()

	f() //如果f函数出现问题，可以从上面的函数中恢复过来
	return
}

func main() {
	fmt.Println(throwsPanic(init1))
	//init1()
}
