package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type IntegrationsCommand struct {
	List   IntegrationsListCommand   `cmd:"list"`
	Get    IntegrationsGetCommand    `cmd:"get"`
	Delete IntegrationsDeleteCommand `cmd:"delete"`
}

type IntegrationsGetCommand struct {
	common.PortalFlags
	Name string `arg:"name"`
}

type IntegrationsListCommand struct {
	common.PortalFlags
}

type IntegrationsDeleteCommand struct {
	common.PortalFlags
	Name string `arg:"name"`
}

func (c *IntegrationsGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetIntegrationWithResponse(context.TODO(), c.Org, c.PortalName, c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Integration)(resp.JSON200))
}

func (c *IntegrationsListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetIntegrationsResponse{}
	integrations := output.IntegrationList{}
	for {
		params := flightdeckv1.GetIntegrationsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetIntegrationsWithResponse(context.TODO(), c.Org, c.PortalName, &params)
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
	resp, err := ctx.APIClient.DeleteIntegrationWithResponse(context.TODO(), c.Org, c.PortalName, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
