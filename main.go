package main

import (
 "github.com/shinyleefeon/gator_rss/internal/config"
 "fmt"
)

func main() {
	
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	cfg.Current_user_name = "Ri"
	err = config.Write(*cfg)
	if err != nil {
		panic(err)
	}

	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Config file: %+v\n", cfg)

}