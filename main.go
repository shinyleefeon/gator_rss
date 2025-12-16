package main

import "github.com/shinyleefeon/rss_gator/internal/config"

func main() {
	readConfig, err := config.Read()
	if err != nil {
		panic(err)
	}
	// This is a placeholder for the main function.
}