package commands

import (
	"os"
	"path"

	"github.com/arctir/flightdeck-cli/commands/common"
	flightdeckclient "github.com/arctir/go-flightdeck/pkg/client"
)

type Cli struct {
	common.Globals

	Auth              AuthCommand              `cmd:"auth"`
	Clusters          ClustersCommand          `cmd:"clusters"`
	Orgs              OrgsCommand              `cmd:"orgs"`
	TenantUsers       TenantUsersCommand       `cmd:"tenantusers"`
	Tenants           TenantsCommand           `cmd:"tenants"`
	Portals           PortalsCommand           `cmd:"portals"`
	Integrations      IntegrationsCommand      `cmd:"integrations"`
	CatalogProviders  CatalogProvidersCommand  `cmd:"catalogproviders"`
	PluginDefinitions PluginDefinitionsCommand `cmd:"plugindefinitions"`
	IdentityProviders IdentityProvidersCommand `cmd:"identityproviders"`
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
