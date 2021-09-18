package main

import (
	"bufio"
	"fmt"
	"os"

	"./calculate"
)

func main() {
	expression, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Print(calculate.Calculate(expression))
}
