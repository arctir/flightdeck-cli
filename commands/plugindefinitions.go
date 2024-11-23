package commands

import (
	"github.com/arctir/flightdeck-cli/commands/common"
)

type PluginDefinitionsListCommand struct {
	common.OrgFlags
}

type PluginDefinitionsGetCommand struct {
	common.OrgFlags
	Name    string `arg:"name"`
	Version int    `arg:"version"`
}

type PluginDefinitionsCommand struct {
	List PluginDefinitionsListCommand `cmd:"list"`
	Get  PluginDefinitionsGetCommand  `cmd:"get"`
}

func (c *PluginDefinitionsListCommand) Run(parent *Cli, ctx *Context) error {
	return nil
	/*
		var err error
		resp := &flightdeckv1.GetPortalVersionPluginDefinitionResponse{}
		definitions := output.PluginDefinitionList{}
		for {
			params := flightdeckv1.GetPortalVersionPluginDefinitionsParams{}
			if resp.JSON200 != nil {
				params.Prev = resp.JSON200.PageInfo.Prev
				params.Next = resp.JSON200.PageInfo.Next
			}
			resp, err = ctx.APIClient.GetPortalVersionPluginDefinitionWithResponse(context.TODO())   GetPortalVersionPluginDefinitionsWithResponse(context.TODO(), c.Org, &params)
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
	*/
}

func (c *PluginDefinitionsGetCommand) Run(parent *Cli, ctx *Context) error {
	return nil
	/*
		resp, err := ctx.APIClient.GetPluginDefinitionWithResponse(context.TODO(), c.Org, c.Name, c.Version)
		if err != nil {
			return err
		}

		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		return output.OutputResult(parent.OutputFormat, (*output.PluginDefinition)(resp.JSON200))
	*/
}
