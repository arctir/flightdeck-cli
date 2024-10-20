package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/google/uuid"
)

type OrgsCreateCommand struct {
	Name      string `arg:"name"`
	ClusterID string `arg:"cluster-id"`
}

type OrgsListCommand struct{}

type OrgsGetCommand struct {
	Id string `arg:"id"`
}

type OrgsDeleteCommand struct {
	Id string `arg:"id"`
}

type OrgsCommand struct {
	Create OrgsCreateCommand `cmd:"create"`
	List   OrgsListCommand   `cmd:"list"`
	Get    OrgsGetCommand    `cmd:"get"`
	Delete OrgsDeleteCommand `cmd:"delete"`
}

func (c *OrgsCreateCommand) Run(parent *Cli, ctx *Context) error {

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

func (c *OrgsListCommand) Run(parent *Cli, ctx *Context) error {
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

func (c *OrgsGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetOrganizationByIDWithResponse(context.TODO(), c.Id)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Organization)(resp.JSON200))
}

func (c *OrgsDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeleteOrganizationByIDWithResponse(context.TODO(), c.Id)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
