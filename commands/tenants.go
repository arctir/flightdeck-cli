package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type TenantsCreateCommand struct {
	common.OrgFlags
	Name        string `arg:"name"`
	DisplayName string `arg:"display_name"`
}

type TenantsGetCommand struct {
	common.OrgFlags
	Name string `arg:"name"`
}

type TenantsListCommand struct {
	common.OrgFlags
}

type TenantsDeleteCommand struct {
	common.OrgFlags
	Name string `arg:"name"`
}

type TenantsCommand struct {
	Create TenantsCreateCommand `cmd:"create"`
	Get    TenantsGetCommand    `cmd:"get"`
	List   TenantsListCommand   `cmd:"list"`
	Delete TenantsDeleteCommand `cmd:"delete"`
}

func (c TenantsCreateCommand) Run(parent *Cli, ctx *Context) error {
	tenantInput := flightdeckv1.TenantInput{
		Name:        c.Name,
		DisplayName: c.DisplayName,
	}

	resp, err := ctx.APIClient.CreateTenantWithResponse(context.TODO(), c.Org, tenantInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Tenant)(resp.JSON201))
}

func (c TenantsGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetTenantWithResponse(context.TODO(), c.Org, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Tenant)(resp.JSON200))
}

func (c TenantsListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetTenantsResponse{}
	tenants := output.TenantList{}
	for {
		params := flightdeckv1.GetTenantsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetTenantsWithResponse(context.TODO(), c.Org, &params)
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
	resp, err := ctx.APIClient.DeleteTenantWithResponse(context.TODO(), c.Org, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
