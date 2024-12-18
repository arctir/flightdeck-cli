package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type TenantUsersCreateCommand struct {
	common.TenantFlags
	Username string `arg:"username"`
	Email    string `arg:"email"`
}

type TenantUsersGetCommand struct {
	common.TenantFlags
	Username *string `arg:"username"`
}

type TenantUsersDeleteCommand struct {
	common.TenantFlags
	Username string `arg:"username"`
}

func (c TenantUsersCreateCommand) Run(parent *Cli, ctx *Context) error {
	user := flightdeckv1.TenantUserInput{
		Username: c.Username,
		Email:    openapi_types.Email(c.Email),
	}

	resp, err := ctx.APIClient.CreateTenantUserWithResponse(context.TODO(), c.Org.String(), c.TenantName, user)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.TenantUser)(resp.JSON201))
}

func (c TenantUsersGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Username == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c TenantUsersGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Username == nil {
		return errors.New("username is required")
	}
	resp, err := ctx.APIClient.GetTenantUserWithResponse(context.TODO(), c.Org.String(), c.TenantName, *c.Username)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.TenantUser)(resp.JSON200))
}

func (c TenantUsersGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetTenantUsersResponse{}
	users := output.TenantUserList{}
	for {
		params := flightdeckv1.GetTenantUsersParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetTenantUsersWithResponse(context.TODO(), c.Org.String(), c.TenantName, &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, t := range *resp.JSON200.Items {
			users = append(users, t)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &users)
}

func (c TenantUsersDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteTenantUserWithResponse(context.TODO(), c.Org.String(), c.TenantName, c.Username)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
