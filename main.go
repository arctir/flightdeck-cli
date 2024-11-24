package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/arctir/flightdeck-cli/commands"
)

const defaultAPIEndpoint = "https://api.arctir.cloud/v1"
const defaultLocalAPIEndpoint = "http://localhost:9090/v1"
const defaultAuthEndpoint = "https://auth.arctir.cloud/realms/arctir-prod"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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

	cli := commands.Cli{}
	commandContext := commands.Context{}
	ctx := kong.Parse(&cli,
		kong.Name("flightdeck"),
		kong.Vars{
			"apiEndpoint":  apiEndpoint,
			"authEndpoint": authEndpoint,
			"configPath":   *configPath,
			"buildVersion": version,
			"buildCommit":  commit,
			"buildDate":    date,
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
