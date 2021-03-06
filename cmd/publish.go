package cmd

import (
	"fmt"
	"github.com/gabrielcolson/c3pm/cli/api"
	"github.com/gabrielcolson/c3pm/cli/config"
	"github.com/gabrielcolson/c3pm/cli/ctpm"
	"net/http"
)

type PublishCmd struct {
	Ignore []string `kong:"optional,name='ignore',short='i',help='Ignore file'"`
}

func (p *PublishCmd) Run() error {
	token, err := config.TokenStrict()
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}
	client := api.API{Client: &http.Client{}, Token: token}
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return ctpm.Publish(pc, client, ctpm.PublishOptions{
		Ignore: p.Ignore,
	})
}
