package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("parseConfig", err.Error())
	}
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	type User struct {
		id       int32
		username string
		emails   []string
	}

	users := []User{}
	rows, _ := dbpool.Query(context.Background(), "select id, username, emails from users")
	for rows.Next() {
		var id int32
		var username string
		var emails []string
		err := rows.Scan(&id, &username, &emails)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		users = append(users, User{id: id, username: username, emails: emails})
	}

	log.Printf("%+v", users)

	// r := gin.Default()
	// r.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "OK")
	// })
	// r.Run(":" + port)
}
