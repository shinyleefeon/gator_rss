package main

import (
 "github.com/shinyleefeon/gator_rss/internal/config"
 "fmt"
 "os"
 "database/sql"
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
	InitCommands.register("reset", deleteUsers)
	InitCommands.register("users", getAllUsers)
	InitCommands.register("agg", aggregateFeeds)
	InitCommands.register("addfeed", middlewareLoggedIn(addFeed))
	InitCommands.register("feeds", listFeeds)
	InitCommands.register("follow", middlewareLoggedIn(followFeed))
	InitCommands.register("following", middlewareLoggedIn(listFollowing))
	InitCommands.register("unfollow", middlewareLoggedIn(unfollowFeed))

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