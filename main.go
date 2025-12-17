package main

import (
 "github.com/shinyleefeon/gator_rss/internal/config"
 "fmt"
 "errors"
 "os"
 "context"
 "time"
 "database/sql"
 "github.com/google/uuid"
 "github.com/shinyleefeon/gator_rss/internal/database"
)

import _ "github.com/lib/pq"

type state struct {
	config *config.Config
	db *database.Queries

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

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("User %s does not exist", cmd.args[0])
		os.Exit(1)
	}


	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username set to:", s.config.Current_user_name)
	return nil
}




func registerUser(s *state, cmd command) error {
	if len(cmd.args) < 1 { 
		return errors.New("Username is required")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("User %s already exists", username)
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: username,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	s.config.SetUser(username)
	fmt.Println("User registered and set as current user:", username)
	fmt.Printf("User details: %+v\n", user)

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

	// Initialize database connection, takes database package and initializes it to state struct
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		panic(err)
	}
	s.db = database.New(db)
	
	
	
	InitCommands := commands{command_map: make(map[string]func(*state, command) error)}
	InitCommands.register("login", handlerLogin)
	InitCommands.register("register", registerUser)

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