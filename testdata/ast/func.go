package ast

import "fmt"

type Cat struct{}

func (c *Cat) Run() {
	fmt.Println("cat is running...")
}

func Hello(msg string) {
	fmt.Println(msg)
}
