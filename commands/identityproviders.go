package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type IdentityProvidersCreateCommand struct {
	common.TenantFlags
	Name string `arg:"name"`
	Type string `arg:"type"`
}

type IdentityProvidersListCommand struct {
	common.TenantFlags
}

type IdentityProvidersGetCommand struct {
	common.TenantFlags
	Name string `arg:"name"`
}

type IdentityProviderDeleteCommand struct {
	common.TenantFlags
	Name string `arg:"name"`
}

type IdentityProvidersCommand struct {
	Create IdentityProvidersCreateCommand `cmd:"create"`
	List   IdentityProvidersListCommand   `cmd:"list"`
	Get    IdentityProvidersGetCommand    `cmd:"get"`
	Delete IdentityProviderDeleteCommand  `cmd:"delete"`
}

func (c *IdentityProvidersCreateCommand) Run(parent *Cli, ctx *Context) error {
	provider := flightdeckv1.IdentityProviderInput{
		Name:           c.Name,
		ProviderConfig: flightdeckv1.IdentityProviderInput_ProviderConfig{},
	}
	resp, err := ctx.APIClient.CreateIdentityProviderWithResponse(context.TODO(), c.Org, c.TenantName, provider)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.IdentityProvider)(resp.JSON201))
}

func (c *IdentityProvidersListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetIdentityProvidersResponse{}
	definitions := output.IdentityProviderList{}
	for {
		params := flightdeckv1.GetIdentityProvidersParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetIdentityProvidersWithResponse(context.TODO(), c.Org, c.TenantName, &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, p := range *resp.JSON200.Items {
			definitions = append(definitions, p)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &definitions)
}

func (c *IdentityProvidersGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetIdentityProviderWithResponse(context.TODO(), c.Org, c.TenantName, c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.IdentityProvider)(resp.JSON200))
}

func (c IdentityProviderDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteIdentityProviderWithResponse(context.TODO(), c.Org, c.TenantName, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
