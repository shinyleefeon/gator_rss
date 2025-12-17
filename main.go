package main

import (
 "github.com/shinyleefeon/gator_rss/internal/config"
 "fmt"
 "errors"
 "os"
)

type state struct {
	config *config.Config

}

type command struct {
	name string
	args []string
}

type commands struct {
	command_map map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 { 
		return errors.New("Username is required")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username set to:", s.config.Current_user_name)
	return nil
}

func (c commands) run(s *state, cmd command) error {
	if handler, exists := c.command_map[cmd.name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c commands) register(name string, f func(*state, command) error) {
	c.command_map[name] = f
}

func main() {
	
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	s := &state{config: cfg}

	InitCommands := commands{command_map: make(map[string]func(*state, command) error)}
	InitCommands.register("login", handlerLogin)

	input := os.Args
	if len(input) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cmd := command{name: input[1], args: input[2:]}

	err = InitCommands.run(s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	

}