package main

import "fmt"

const (
	WHITE = iota //0
	BLACK
	RED
	YELLOW
	GREEN //4
)

type Color byte

type Box struct {
	Len, wid, height float32
	c                Color
}

//define your own data type
type BoxList []Box //Slice

//set the function's receiver
func (b Box) Volume() float32 {
	return b.Len * b.wid * b.height
}

func (b *Box) SetColor(co Color) {
	//Attention: when you write *b.c this way, the compiler will send error.
	(*b).c = co
}

//get the volume biggest's color
func (bL BoxList) BiggestsColor() (c Color) {
	var v float32
	v = 0.0
	c = Color(WHITE)

	for _, b := range bL {
		if v < b.Volume() {
			c = b.c
			v = b.Volume()
		}
	}

	return c

}

func (bL BoxList) PrintItBlack() {
	for i, _ := range bL {
		bL[i].SetColor(BLACK)
	}
}

//any datatype ,you can give them a reciver ,so that that data type can do something more
func (c Color) String() string {
	str := []string{"WHITE", "BLACK", "RED", "YELLOW", "GREEN"}
	return str[c]
}

func main() {
	boxes := BoxList{
		Box{10, 2, 3, WHITE},
		Box{4, 2, 5, BLACK},
		Box{5, 3, 6, RED},
		Box{6, 2, 1, YELLOW},
		Box{8, 3, 1, GREEN},
	}

	//attention append this function's usage.
	boxes = append(boxes, Box{2, 3, 4, GREEN})

	fmt.Printf("now we have %d boxes\n", len(boxes))
	fmt.Println("the first box volume is ", boxes[0].Volume())
	fmt.Println("the biggest box's color is ", boxes.BiggestsColor())
	fmt.Println("the last box's color is ", boxes[len(boxes)-1].c.String())
	fmt.Println("paint the all box to black")
	boxes.PrintItBlack()
	//array,map and so on we can print them directly
	fmt.Println(boxes)
}
