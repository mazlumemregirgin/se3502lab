package main

import (
	"fmt"
	"se3502lab/week3/lab"
)

func main() {

	todos := []string{"server side", "big data", "vbt", "kedi kumu alınacak", "limon tuzu alınacak"}

	fmt.Println(todos)

	todos = append(todos[:2], todos[3:]...)

	fmt.Println(todos)

	//lab.Capacity()
	lab.OrderMap()

}
