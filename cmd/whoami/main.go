package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/authn/pkg/oauth2"
)

type envConfig struct {
	FilePath string `envconfig:"FILE_PATH" default:"/var/run/ko/" required:"true"`
	AppPort  string `envconfig:"APP_PORT" default:"8181"`
}

func main() {
	ctx := context.Background()
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	staticDir := http.Dir(path.Join(env.FilePath, "/static"))
	templates := template.Must(template.ParseGlob(filepath.Join(env.FilePath, "/templates", "*")))

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		var output bytes.Buffer
		if err := templates.ExecuteTemplate(&output, "main.html", map[string]interface{}{
			"User": oauth2.AuthenticatedUser(r.Context()),
		}); err != nil {
			http.Error(w, "could not render template", http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(output.Bytes()); err != nil {
			log.Printf("error: %v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	log.Fatal(http.ListenAndServe(":"+env.AppPort, oauth2.NewOrDie(ctx, mux)))
}
