package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shinyleefeon/gator_rss/internal/database"
)

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
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
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

func deleteUsers(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("Error deleting users:", err)
		os.Exit(1)
	}
	fmt.Println("All users deleted from the database.")
	return nil
}

func getAllUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Registered users:")
	for _, user := range users {
		if user == s.config.Current_user_name {
			fmt.Println("* ", user, "(current)")
		} else {
			fmt.Println("* ", user)
		}
	}
	return nil
}
