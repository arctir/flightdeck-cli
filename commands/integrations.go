package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type IntegrationsCreateCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the integration."`
	Type string `arg:"type" help:"Type of the integration."`
}
type IntegrationsGetCommand struct {
	common.PortalFlags
	Name *string `arg:"name" optional:"" help:"Name of the integration to get. If not provided, lists all integrations."`
}

type IntegrationsDeleteCommand struct {
	common.PortalFlags
	Name string `arg:"name" help:"Name of the integration."`
}

func (c IntegrationsCreateCommand) Run(ctx *Context) error {
	return errors.New("not implemented")
}

func (c IntegrationsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c IntegrationsGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return errors.New("name is required")
	}
	resp, err := ctx.APIClient.GetIntegrationWithResponse(context.TODO(), c.Org.String(), c.PortalName, *c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Integration)(resp.JSON200))
}

func (c IntegrationsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetIntegrationsResponse{}
	integrations := output.IntegrationList{}
	for {
		params := flightdeckv1.GetIntegrationsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetIntegrationsWithResponse(context.TODO(), c.Org.String(), c.PortalName, &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, i := range *resp.JSON200.Items {
			integrations = append(integrations, i)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &integrations)
}

func (c IntegrationsDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteIntegrationWithResponse(context.TODO(), c.Org.String(), c.PortalName, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
