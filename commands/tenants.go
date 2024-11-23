package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type TenantsCreateCommand struct {
	common.OrgFlags
	Name        string `arg:"name" help:"Name of the tenant."`
	DisplayName string `arg:"display_name" help:"Display name of the tenant."`
}

type TenantsGetCommand struct {
	common.OrgFlags
	Name *string `arg:"name" optional:"" help:"Name of the tenant."`
}

type TenantsDeleteCommand struct {
	common.OrgFlags
	Name string `arg:"name" help:"Name of the tenant."`
}

func (c TenantsCreateCommand) Run(parent *Cli, ctx *Context) error {
	tenantInput := flightdeckv1.TenantInput{
		Name:        c.Name,
		DisplayName: c.DisplayName,
	}

	resp, err := ctx.APIClient.CreateTenantWithResponse(context.TODO(), c.Org.String(), tenantInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Tenant)(resp.JSON201))
}

func (c TenantsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c TenantsGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return errors.New("name is required")
	}
	resp, err := ctx.APIClient.GetTenantWithResponse(context.TODO(), c.Org.String(), *c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Tenant)(resp.JSON200))
}

func (c TenantsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetTenantsResponse{}
	tenants := output.TenantList{}
	for {
		params := flightdeckv1.GetTenantsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetTenantsWithResponse(context.TODO(), c.Org.String(), &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, t := range *resp.JSON200.Items {
			tenants = append(tenants, t)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &tenants)
}

func (c TenantsDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteTenantWithResponse(context.TODO(), c.Org.String(), c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
