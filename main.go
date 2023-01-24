package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akerl/prospectus/v3/plugin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type config struct {
	Seconds int    `json:"seconds"`
	Branch  string `json:"branch"`
}

func run(cfg config) error {
	if cfg.Seconds == 0 {
		return fmt.Errorf("required arg not provided: seconds")
	}

	branch := cfg.Branch
	if branch == "" {
		branch = "main"
	}
	refName := fmt.Sprintf("refs/heads/%s", branch)

	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	ref, err := repo.Reference(plumbing.NewBranchReferenceName(refName), false)
	if err != nil {
		return err
	}

	log, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}
	defer log.Close()

	last, err := log.Next()
	if err != nil {
		return err
	}

	stamp := last.Committer.When

	currentTime := time.Now()
	age := int(currentTime.Sub(stamp).Seconds())

	if age <= cfg.Seconds {
		fmt.Println("fresh")
	} else {
		fmt.Println("outdated")
	}
	return nil
}

func main() {
	cfg := config{}
	err := plugin.ParseConfig(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s", err)
		os.Exit(1)
	}
	err = run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
