package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type CatalogProvidersCommand struct {
	List CatalogProvidersListCommand `cmd:"list"`
	Get  CatalogProvidersGetCommand  `cmd:"update"`
}

type CatalogProvidersListCommand struct {
	common.PortalFlags
}

type CatalogProvidersGetCommand struct {
	common.PortalFlags
	Name string `arg:"name"`
}

func (c *CatalogProvidersListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetCatalogProvidersResponse{}
	catalogProviders := output.CatalogProviderList{}
	for {
		params := flightdeckv1.GetCatalogProvidersParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetCatalogProvidersWithResponse(context.TODO(), c.Org, c.PortalName, &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, i := range *resp.JSON200.Items {
			catalogProviders = append(catalogProviders, i)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &catalogProviders)
}

func (c *CatalogProvidersGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetCatalogProviderWithResponse(context.TODO(), c.Org, c.PortalName, c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.CatalogProvider)(resp.JSON200))
}
