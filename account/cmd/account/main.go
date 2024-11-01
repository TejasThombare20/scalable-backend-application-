package main

import (
	"log"
	"time"

	"github.com/TejasThombare20/go-microservice/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	DatabaseURL := "postgres://postgres:Tejas@account_db:5432/go-microservices?sslmode=disable"

	var r account.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {

		log.Println("cfg.DatabasURL", cfg.DatabaseURL)

		r, err = account.NewPostgresRepository(DatabaseURL)

		if err != nil {
			log.Println("err while connecting with account db postgres", err)
		}
		return
	})
	defer r.Close()
	log.Println("Listing on port 8080")
	s := account.NewService(r)

	log.Fatal(account.ListernGRPC(s, 8080))
}
