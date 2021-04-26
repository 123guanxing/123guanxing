package main

import "fmt"

func main() {
	i := 0
Loop:
	for {
		fmt.Println(i)
		i++
		if i >= 10 {
			//continue Loop
			break Loop
		}
	}
}
