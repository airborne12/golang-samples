package main

import (
	"fmt"
	"github.com/airborne12/golang-samples/configmodule/config"
)

func main(){
	config := config.Get()
	fmt.Println(config)
}