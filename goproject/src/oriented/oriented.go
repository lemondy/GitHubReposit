// oriented project oriented.go
package oriented

import "fmt"

//define your own type
type months map[string]int

//in go there is not enum type, while we can use const and itoa keyword to define const

const (
	Sunday = iota
	Monday
	Tuseday
	Wednesday
	Thursday
	Friday
	Saturday
)

func Main() {
	m := months{
		"Januray":   31,
		"Feburary":  28,
		"March":     31,
		"April":     30,
		"May":       31,
		"June":      30,
		"July":      31,
		"August":    31,
		"September": 30,
		"October":   31,
		"November":  30,
		"December":  31,
	}
	//add element
	m["other"] = 29
	for key, month := range m {
		fmt.Println(key, month)
	}

	fmt.Println(m)
	delete(m, "other")
	/*for key, month := range m {
		fmt.Println(key, month)
	}*/
	fmt.Println(m)

}
