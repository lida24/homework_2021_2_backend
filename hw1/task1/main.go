package main

import (
	"fmt"
	"os"

	"./app"
)

func main() {
	application := app.Application{}
	fmt.Println(application.RunApp(os.Stdin))
}
