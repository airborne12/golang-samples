package main

import (
	"fmt"

	"github.com/airborne12/golang-samples/configmodule/config"
)

var gConfig = config.Get()

func main() {
	fmt.Println("listen address:", gConfig.ListenAddr)
}
