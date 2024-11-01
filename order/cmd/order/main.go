package main

import (
	"log"
	"time"

	"github.com/TejasThombare20/go-microservice/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envConfig:"DATABASE_URL"`
	AccountURL  string `envConfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envConfig:"CATALOG_SERVICE_URL"`
}

func main() {

	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository

	DatabaseURL := "postgres://postgres:Tejas@order_db:5432/go-microservices?sslmode=disable"

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {

		log.Println("cfg.DatabaseURL", cfg.DatabaseURL)

		r, err = order.NewPostgresRepository(DatabaseURL)

		if err != nil {
			log.Println("error while connecting with order postgres db ", err)
		}
		return err
	})

	defer r.Close()
	log.Println("Listening on port 8080... ")

	s := order.NewService(r)

	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
