package main

import (
	"database/sql"
	"fmt"
	"github.com/JakeBurrell/gator/internal/config"
	"github.com/JakeBurrell/gator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	gCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	gCommands.register("login", handlerLogin)
	gCommands.register("register", handlerRegister)
	gCommands.register("reset", handlerRest)
	gCommands.register("users", handlerUsers)
	gCommands.register("agg", handlerAgg)
	gCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	gCommands.register("feeds", handleFeeds)
	gCommands.register("follow", middlewareLoggedIn(handleFollow))
	gCommands.register("following", middlewareLoggedIn(handlerFollowing))
	gCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	db, err := sql.Open("postgres", cfg.DataBaseURL)
	if err != nil {
		log.Fatalf("Error connection to the databse: %v", err)
	}
	dbQueries := database.New(db)
	gState := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("No commands provided")
		os.Exit(1)
	}

	err = gCommands.run(
		&gState,
		command{
			name: args[1],
			args: args[2:],
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
