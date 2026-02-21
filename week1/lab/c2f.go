package lab

import "fmt"

func C2F() {
	var celsius float64
	fmt.Print("Enter temperature in Celsius: ")
	fmt.Scanln(&celsius)
	fahrenheit := (celsius * 9 / 5) + 32
	fmt.Printf("%.2f Celsius is %.2f Fahrenheit\n", celsius, fahrenheit)

}
