package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/environ"
	"os"
)

const (
	default_endpoint = "http://127.0.0.1:8080/envs"
	default_secret   = "my_secret"
)

func main() {
	endpoint := flag.String("endpoint", os.Getenv("ENDPOINT"), "endpoint")
	secret := flag.String("secret", os.Getenv("EXT_ENV_SECRET_KEY"), "secret")
	skipverify := flag.Bool("skipverify", false, "skipverify")

	repoSlug := flag.String("repo-slug", os.Getenv("REPO_SLUG"), "repo-slug")

	flag.Parse()

	if *endpoint == "" {
		*endpoint = default_endpoint
	}
	if *secret == "" {
		*secret = default_secret
	}

	client := environ.Client(*endpoint, *secret, *skipverify)

	ctx := context.Background()
	r := &environ.Request{}
	r.Repo = drone.Repo{Slug: *repoSlug}

	list, err := client.List(ctx, r)
	if err != nil {
		fmt.Printf("err:  %v", err)
	}
	for _, variable := range list {
		fmt.Printf("Name: %s\nData:\n%s\nMask: %v\n---\n", variable.Name, variable.Data, variable.Mask)
	}
}
