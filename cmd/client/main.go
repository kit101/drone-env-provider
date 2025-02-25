package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/environ"
	"gopkg.in/yaml.v3"
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

	format := flag.String("format", "line", "format: support line, yaml, json. (default: line)")
	pretty := flag.Bool("pretty", true, "pretty in format json (default: true)")

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
	if *format == "line" {
		printLine(list)
	} else if *format == "yaml" {
		printYaml(list)
	} else if *format == "json" {
		printJson(list, *pretty)
	}
}

func printLine(list []*environ.Variable) {
	for _, variable := range list {
		fmt.Printf("Name: %s\nData:\n%s\nMask: %v\n---\n", variable.Name, variable.Data, variable.Mask)
	}
}

func printYaml(list []*environ.Variable) {
	var node yaml.Node
	b, _ := yaml.Marshal(list)
	_ = yaml.Unmarshal(b, &node)
	plain, err := yaml.Marshal(&node)
	if err != nil {
		fmt.Printf("err:  %v, will do printLine.", err)
		printLine(list)
		return
	}
	fmt.Printf("%s", plain)
}

func printJson(list []*environ.Variable, pretty bool) {
	var plain []byte
	var err error
	if pretty {
		plain, err = json.MarshalIndent(list, "", "  ")
	} else {
		plain, err = json.Marshal(list)
	}
	if err != nil {
		fmt.Printf("err:  %v, will do printLine.", err)
		printLine(list)
	} else {
		fmt.Printf("%s", plain)
	}
}
