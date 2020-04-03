package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/authn/pkg/oauth2"
)

type envConfig struct {
	AppPort string `envconfig:"APP_PORT" default:"8181"`
}

func main() {
	ctx := context.Background()
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := oauth2.AuthenticatedUser(r.Context())
		_, _ = w.Write([]byte(fmt.Sprintf("Hello %s (%s)!", user.Name, user.Email)))
	})

	log.Fatal(http.ListenAndServe(":"+env.AppPort, oauth2.NewOrDie(ctx, mux)))
}
