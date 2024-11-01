package main

import (
	"log"
	"time"

	"github.com/TejasThombare20/go-microservice/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `json:"DATABASE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		log.Println("cfg.DatabasURL", cfg.DatabaseURL)

		r, err = catalog.NewEsasticRepository(cfg.DatabaseURL)

		if err != nil {
			log.Println("error creating elastic repository")
			log.Println(err)
		}

		return
	})

	defer r.Close()
	log.Println("Listening on port 8080 .....")

	s := catalog.NewService(r)

	log.Fatal(catalog.ListerGRPC(s, 8080))

}
