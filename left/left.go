package main

import "fmt"

const (
	one int = 1 << iota
	two
	three
	four
)

func main() {
	b := T()
	if (b & four) == four {
		fmt.Println("true")
		return
	}
	fmt.Println("false")
}

func T () int {
	return 9
}





