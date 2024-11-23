package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type CatalogProvidersCreateCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the catalog provider."`
}

type CatalogProvidersGetCommand struct {
	common.PortalFlags
	Name *string `arg:"name" optional:"" help:"Name of the catalog provider to get. If not provided, lists all catalog providers."`
}

type CatalogProvidersDeleteCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the catalog provider."`
}

func (c CatalogProvidersCreateCommand) Run(parent *Cli, ctx *Context) error {
	return errors.New("not implemented")
}

func (c CatalogProvidersGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c CatalogProvidersGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return errors.New("name is required")
	}
	resp, err := ctx.APIClient.GetCatalogProviderWithResponse(context.TODO(), c.Org.String(), c.PortalName, *c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.CatalogProvider)(resp.JSON200))
}

func (c CatalogProvidersGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetCatalogProvidersResponse{}
	catalogProviders := output.CatalogProviderList{}
	for {
		params := flightdeckv1.GetCatalogProvidersParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetCatalogProvidersWithResponse(context.TODO(), c.Org.String(), c.PortalName, &params)
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

func (c CatalogProvidersDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteCatalogProviderWithResponse(context.TODO(), c.Org.String(), c.PortalName, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
