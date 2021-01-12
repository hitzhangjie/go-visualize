package ast

import "fmt"

func TestIfElse() {

	num := 1
	if num == 1 {
		fmt.Println("1")
	} else if num == 2 {
		fmt.Println("2")
	} else {
		fmt.Println("others")
	}
}
