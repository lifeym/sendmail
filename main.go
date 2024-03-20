package main

import (
	"os"

	"github.com/lifeym/she/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		//fmt.Fprintln(os.Stderr, err)
		// fmt.Println("exit 1")
		os.Exit(1)
		return
	}

	//fmt.Println("exit 0")
	os.Exit(0)
}
