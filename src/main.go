package main

import (
	"fmt"
	"github.com/Juno-chat-app/user-service/config"
)

func main() {
	conf := config.LoadConfiguration("test.txt")
	fmt.Println(conf)
}
