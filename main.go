package main

import (
	"beats_refactor/config"
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if _, err := config.InitConfig(); err != nil {
		fmt.Println(err)
	}
}
