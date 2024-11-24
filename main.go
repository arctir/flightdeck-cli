package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/arctir/flightdeck-cli/commands"
	"github.com/google/uuid"
)

const defaultAPIEndpoint = "https://api.arctir.cloud/v1"
const defaultLocalAPIEndpoint = "http://localhost:9090/v1"
const defaultAuthEndpoint = "https://auth.arctir.cloud/realms/arctir-prod"

func main() {
	configPath, err := commands.ConfigPath()
	if err != nil {
		panic(err)
	}

	apiEndpoint := defaultAPIEndpoint
	authEndpoint := defaultAuthEndpoint

	flightdeckEnv := os.Getenv("FLIGHTDECK_ENV")
	if flightdeckEnv == "local" {
		apiEndpoint = defaultLocalAPIEndpoint
	}

	defaultOrg := os.Getenv("FLIGHTDECK_ORG")
	if defaultOrg == "" {
		defaultOrg = uuid.Nil.String()
	}

	cli := commands.Cli{}
	commandContext := commands.Context{}
	ctx := kong.Parse(&cli,
		kong.Name("flightdeck"),
		kong.Vars{
			"apiEndpoint":  apiEndpoint,
			"authEndpoint": authEndpoint,
			"configPath":   *configPath,
			"defaultOrg":   defaultOrg,
		},
		kong.Bind(&commandContext),
		kong.Bind(&cli.Globals),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			Summary:             false,
			NoExpandSubcommands: true,
		}))
	err = ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
