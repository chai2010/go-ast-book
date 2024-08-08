package main

import "fmt"

const Pi = 3.14

func main() {
	for i := 2; i <= 8; i++ {
		fmt.Printf("%d*Pi = %.2f\n", i, Pi*float64(i))
	}
}
