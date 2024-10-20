package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type PortalsCreateCommand struct {
	common.OrgFlags
	Name             string `arg:"name"`
	Domain           string `arg:"domain"`
	Title            string `arg:"title"`
	OrganizationName string `arg:"org-name"`
	Version          string `arg:"version"`
}

type PortalsListCommand struct {
	common.OrgFlags
}

type PortalsGetCommand struct {
	common.OrgFlags
	Name string `arg:"name"`
}

type PortalDeleteCommand struct {
	common.OrgFlags
	Name string `arg:"name"`
}

type PortalsCommand struct {
	Create PortalsCreateCommand `cmd:"create"`
	List   PortalsListCommand   `cmd:"list"`
	Get    PortalsGetCommand    `cmd:"get"`
	Delete PortalDeleteCommand  `cmd:"delete"`
}

func (c *PortalsCreateCommand) Run(parent *Cli, ctx *Context) error {
	portalInput := flightdeckv1.PortalInput{
		Name:             c.Name,
		Domain:           c.Domain,
		Title:            c.Title,
		OrganizationName: c.OrganizationName,
		Version:          c.Version,
	}
	resp, err := ctx.APIClient.CreatePortalWithResponse(context.TODO(), c.Org, portalInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Portal)(resp.JSON201))
}

func (c *PortalsListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetPortalsResponse{}
	portals := output.PortalList{}
	for {
		params := flightdeckv1.GetPortalsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetPortalsWithResponse(context.TODO(), c.Org, &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, p := range *resp.JSON200.Items {
			portals = append(portals, p)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, &portals)
}

func (c *PortalsGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetPortalWithResponse(context.TODO(), c.Org, c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Portal)(resp.JSON200))
}

func (c PortalDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeletePortalWithResponse(context.TODO(), c.Org, c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
