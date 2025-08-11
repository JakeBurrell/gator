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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%+v\n", feed)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %s <name> <url>\n", cmd.name)
	}

	// Get user id
	usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	// Add Feed
	new_feed := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    usr.ID,
	}

	_, err = s.db.AddFeed(
		context.Background(),
		new_feed,
	)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    usr.ID,
			FeedID:    new_feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't follow new feed: %w", err)
	}

	return nil
}

func handleFollows(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.name)
	}

	//  Gets user id
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	// Gets feed id
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	follows, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("====================\n")
	fmt.Printf("User %s followed feed %s\n", follows.UserName, follows.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	// Gets user for id
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed_followers: %w", err)
	}
	fmt.Printf("You: %s are following\n", user.Name)
	for i, feed := range follows {
		fmt.Printf("	%d: %s ", i+1, feed.Feedname)
	}
	return nil
}
