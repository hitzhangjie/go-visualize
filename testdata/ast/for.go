package ast

import "fmt"

func TestForStmt() {
	for i := 0; i < 10; i++ {
		fmt.Println("hello world")
	}
}
