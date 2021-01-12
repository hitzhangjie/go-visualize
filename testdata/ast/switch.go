package ast

import "fmt"

func TestSwitchCase() {
	num := 1

	switch num {
	case 1:
		fmt.Println("1")
	case 2:
		fmt.Println("2")
	default:
		fmt.Println("others")
	}
}
