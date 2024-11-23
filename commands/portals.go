package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/confluentinc/go-editor"
	"sigs.k8s.io/yaml"
)

type PortalsCreateCommand struct {
	common.OrgFlags
	Name             string `arg:"name"`
	Domain           string `arg:"domain"`
	Title            string `arg:"title"`
	OrganizationName string `arg:"org-name"`
	VersionID        string `arg:"version-id"`
}

type PortalsGetCommand struct {
	common.OrgFlags
	Name *string `arg:"name" optional:"" help:"Name of the portal to get. If not provided, lists all portals."`
}

type PortalDeleteCommand struct {
	common.OrgFlags
	Name string `arg:"name" help:"Name of the portal to delete."`
}

func (c PortalsCreateCommand) Run(parent *Cli, ctx *Context) error {

	out, err := yaml.Marshal(flightdeckv1.PortalInput{})
	original := bytes.NewBufferString(string(out))

	edit := editor.NewEditor()
	edited, file, err := edit.LaunchTempFile("example", original)
	defer os.Remove(file)
	if err != nil {
		return err
	}
	fmt.Println("edited:", string(edited))
	return nil

	portalInput := flightdeckv1.PortalInput{
		Name:             c.Name,
		Domain:           c.Domain,
		Title:            c.Title,
		OrganizationName: c.OrganizationName,
		VersionId:        c.VersionID,
	}
	resp, err := ctx.APIClient.CreatePortalWithResponse(context.TODO(), c.Org.String(), portalInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Portal)(resp.JSON201))
}

func (c PortalsGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Name == nil {
		return c.list(parent, ctx)
	}
	return c.get(ctx, parent)
}

func (c PortalsGetCommand) list(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetPortalsResponse{}
	portals := output.PortalList{}
	for {
		params := flightdeckv1.GetPortalsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetPortalsWithResponse(context.TODO(), c.Org.String(), &params)
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

func (c PortalsGetCommand) get(ctx *Context, parent *Cli) error {
	if c.Name == nil {
		return errors.New("name is required")
	}
	resp, err := ctx.APIClient.GetPortalWithResponse(context.TODO(), c.Org.String(), *c.Name)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Portal)(resp.JSON200))
}

func (c PortalDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeletePortalWithResponse(context.TODO(), c.Org.String(), c.Name)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
