package main

import (
	"context"
	"fmt"
	"github.com/JakeBurrell/gator/internal/config"
	"github.com/JakeBurrell/gator/internal/database"
	"github.com/JakeBurrell/gator/internal/rss"
	"github.com/google/uuid"
	"time"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Failed to log in user database return: %v\n", err)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("The user: %s was logged in\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("User already exists\n")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	_, err = s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	return nil

}

func handlerRest(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Users reset failed with: %v", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	cfg, err := config.Read()
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Print(user)
		if user == cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rss, err := rss.FetchFeed(
		context.Background(),
		"https://www.wagslane.dev/index.xml",
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", rss)
	return nil

}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.name)
	}

	// Get current username
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	username := cfg.CurrentUserName

	// Get user id
	usr, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	// Add Feed
	_, err = s.db.AddFeed(
		context.Background(),
		database.AddFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    usr.ID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%+v\n", feed)
	}

	return nil
}
