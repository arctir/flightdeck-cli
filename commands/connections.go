package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type ConnectionsCreateCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the connection."`
	Type string `arg:"type" help:"Type of the connection."`
}

type ConnectionsGetCommand struct {
	common.PortalFlags
	Name *string `arg:"name" optional:"" help:"Name of the connection to get. If not provided, lists all connections."`
}

type ConnectionsDeleteCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the connection."`
}

func (c ConnectionsCreateCommand) Run(parent *Cli, ctx *Context) error {
	return errors.New("not implemented")
}

func (c ConnectionsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c ConnectionsGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return errors.New("name is required")
	}
	resp, err := ctx.APIClient.GetConnectionWithResponse(context.TODO(), c.Org.String(), c.PortalName, *c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Connection)(resp.JSON200))
}

func (c ConnectionsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetConnectionsResponse{}
	catalogProviders := output.ConnectionList{}
	for {
		params := flightdeckv1.GetConnectionsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetConnectionsWithResponse(context.TODO(), c.Org.String(), c.PortalName, &params)
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

func (c ConnectionsDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteConnectionWithResponse(context.TODO(), c.Org.String(), c.PortalName, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
