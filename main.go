package main

import (
	"fmt"
	"os"
	"github.com/jc-m/rhub/cmd"
)


func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
