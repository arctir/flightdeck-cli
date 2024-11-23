package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type PortalVersionsGetCommand struct {
	Id   *string `arg:"id" optional:"" help:"ID of the portal version to get. If not provided, lists all portal versions."`
	Name *string `arg:"name" optional:"" help:"Name of the portal version to get. If not provided, lists all portal versions."`
}

func (c PortalVersionsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Id == nil {
		return c.list(parent, ctx)
	}
	return c.get(ctx, parent)
}

func (c PortalVersionsGetCommand) get(ctx *Context, parent *Cli) error {
	if c.Id == nil {
		return errors.New("id is required")
	}
	resp, err := ctx.APIClient.GetPortalVersionWithResponse(context.TODO(), *c.Id)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.PortalVersion)(resp.JSON200))
}

func (c PortalVersionsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetPortalVersionsResponse{}
	versions := output.PortalVersionList{}

	var versionName *string
	if c.Name == nil {
		versionName = c.Name
	}

	for {
		params := flightdeckv1.GetPortalVersionsParams{
			VersionName: versionName,
		}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetPortalVersionsWithResponse(context.TODO(), &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, o := range *resp.JSON200.Items {
			versions = append(versions, o)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &versions)
}
