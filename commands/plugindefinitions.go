package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/arctir/flightdeck-cli/commands/common"
	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"sigs.k8s.io/yaml"
)

type PluginDefinitionsCreateCommand struct {
	common.OrgFlags
	Filename string `arg:"filename"`
}

type PluginDefinitionsListCommand struct {
	common.OrgFlags
}

type PluginDefinitionsGetCommand struct {
	common.OrgFlags
	Name    string `arg:"name"`
	Version int    `arg:"version"`
}

type PluginDefinitionsUpdateCommand struct {
	common.OrgFlags
	Filename string `arg:"filename"`
	Name     string `arg:"name"`
	Version  int    `arg:"version"`
}

type PluginDefinitionDeleteCommand struct {
	common.OrgFlags
	Name    string `arg:"name"`
	Version int    `arg:"version"`
}

type PluginDefinitionsCommand struct {
	Create PluginDefinitionsCreateCommand `cmd:"create"`
	List   PluginDefinitionsListCommand   `cmd:"list"`
	Get    PluginDefinitionsGetCommand    `cmd:"get"`
	Update PluginDefinitionsUpdateCommand `cmd:"update"`
	Delete PluginDefinitionDeleteCommand  `cmd:"delete"`
}

func readYamlInput(file string, target interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, target)
	if err != nil {
		return err
	}

	out, _ := json.MarshalIndent(target, "", "    ")
	fmt.Println(string(out))
	return err
}

func (c *PluginDefinitionsCreateCommand) Run(parent *Cli, ctx *Context) error {
	definitionInput := flightdeckv1.PluginDefinitionInput{}

	err := readYamlInput(c.Filename, &definitionInput)
	if err != nil {
		return err
	}

	resp, err := ctx.APIClient.CreatePluginDefinitionWithResponse(context.TODO(), c.Org, definitionInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.PluginDefinition)(resp.JSON201))
}

func (c *PluginDefinitionsListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetPluginDefinitionsResponse{}
	definitions := output.PluginDefinitionList{}
	for {
		params := flightdeckv1.GetPluginDefinitionsParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetPluginDefinitionsWithResponse(context.TODO(), c.Org, &params)
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

func (c *PluginDefinitionsGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetPluginDefinitionWithResponse(context.TODO(), c.Org, c.Name, c.Version)
	if err != nil {
		return err
	}

	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.PluginDefinition)(resp.JSON200))
}

func (c *PluginDefinitionsUpdateCommand) Run(parent *Cli, ctx *Context) error {
	definitionInput := flightdeckv1.PluginDefinitionInput{}

	err := readYamlInput(c.Filename, &definitionInput)
	if err != nil {
		return err
	}

	resp, err := ctx.APIClient.UpdatePluginDefinitionWithResponse(context.TODO(), c.Org, c.Name, c.Version, definitionInput)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.PluginDefinition)(resp.JSON200))
}

func (c PluginDefinitionDeleteCommand) Run(ctx *Context) error {
	resp, err := ctx.APIClient.DeletePluginDefinitionWithResponse(context.TODO(), c.Org, c.Name, c.Version)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return nil
}
