package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/alecthomas/kong"
	"github.com/arctir/flightdeck-cli/commands/common"
	flightdeckclient "github.com/arctir/go-flightdeck/pkg/client"
)

type Cli struct {
	common.Globals

	Auth    AuthCommand    `cmd:"auth"`
	Create  CreateCommand  `cmd:"create"`
	Get     GetCommand     `cmd:"get" predictor:"getPredictor"`
	Delete  DeleteCommand  `cmd:"delete"`
	Version VersionCommand `cmd:"version"`
}

type VersionCommand struct{}

func (c VersionCommand) Run(vars kong.Vars) error {
	fmt.Printf("flightdeck %s, commit %s, built at %s\n", vars["buildVersion"], vars["buildCommit"], vars["buildDate"])
	return nil
}

func (c *Cli) AfterApply(ctx *Context, globals *common.Globals) error {
	if !ctx.SkipConfigCheck {
		config, err := flightdeckclient.ReadConfig(globals.ConfigPath)
		if err != nil {
			return err
		}
		ctx.Config = config

		client, err := flightdeckclient.NewClient(globals.APIEndpoint, *config)
		if err != nil {
			return err
		}
		ctx.APIClient = client
	}
	return nil
}

func ConfigPath() (*string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	flightdeckConfigDir := path.Join(dirname, ".flightdeck")
	err = os.MkdirAll(flightdeckConfigDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	flightdeckConfig := path.Join(flightdeckConfigDir, "config.yaml")
	return &flightdeckConfig, nil
}
