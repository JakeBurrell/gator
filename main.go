package main

import (
	"fmt"

	"github.com/JakeBurrell/gator/internal/config"
)

func main() {
	_, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	err = config.SetUser("Jake")
	if err != nil {
		fmt.Println(err)
	}
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)

}
