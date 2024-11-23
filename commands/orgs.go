package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/google/uuid"
)

type OrgsCreateCommand struct {
	Name      string `arg:"name" help:"Name of the organization."`
	ClusterID string `arg:"cluster-id" help:"ID of the cluster to create the organization in."`
}

type OrgsGetCommand struct {
	Id *uuid.UUID `arg:"id" optional:"" name:"id" help:"ID of the organization to get. If not provided, lists all organizations."`
}

type OrgsDeleteCommand struct {
	Id uuid.UUID `arg:"id" help:"ID of the organization to delete."`
}

func (c OrgsCreateCommand) Run(parent *Cli, ctx *Context) error {
	clusterId, err := uuid.Parse(c.ClusterID)
	if err != nil {
		return err
	}

	orgInput := flightdeckv1.OrganizationInput{
		Name:      c.Name,
		ClusterId: clusterId,
	}

	resp, err := ctx.APIClient.CreateOrganizationWithResponse(context.TODO(), orgInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Organization)(resp.JSON201))
}

func (c OrgsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Id == nil {
		return c.list(parent, ctx)
	}
	return c.get(ctx, parent)
}

func (c OrgsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetOrganizationsResponse{}
	orgs := output.OrganizationList{}
	for {
		params := flightdeckv1.GetOrganizationsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetOrganizationsWithResponse(context.TODO(), &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, o := range *resp.JSON200.Items {
			orgs = append(orgs, o)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &orgs)
}

func (c OrgsGetCommand) get(ctx *Context, parent *Cli) error {
	if c.Id == nil {
		return errors.New("id is required")
	}
	resp, err := ctx.APIClient.GetOrganizationByIDWithResponse(context.TODO(), c.Id.String())
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Organization)(resp.JSON200))
}

func (c *OrgsDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteOrganizationByIDWithResponse(context.TODO(), c.Id.String())
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
