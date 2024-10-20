package commands

import (
	"context"

	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

type ClustersListCommand struct{}

type ClustersGetCommand struct {
	Id string `arg:"id"`
}

type ClustersCommand struct {
	List ClustersListCommand `cmd:"list"`
	Get  ClustersGetCommand  `cmd:"get"`
}

func (c *ClustersListCommand) Run(parent *Cli, ctx *Context) error {
	var err error
	resp := &flightdeckv1.GetClustersResponse{}
	clusters := []flightdeckv1.Cluster{}
	for {
		params := flightdeckv1.GetClustersParams{}
		if resp.JSON200 != nil {
			params.Prev = resp.JSON200.PageInfo.Prev
			params.Next = resp.JSON200.PageInfo.Next
		}
		resp, err = ctx.APIClient.GetClustersWithResponse(context.TODO(), &params)
		if err != nil {
			return err
		}
		if !output.IsOutputtableResponse(resp.HTTPResponse) {
			return nil
		}
		for _, o := range *resp.JSON200.Items {
			clusters = append(clusters, o)
		}
		if resp.JSON200.PageInfo.Next == nil {
			break
		}
	}

	return output.OutputResult(parent.OutputFormat, (*output.ClusterList)(&clusters))
}

func (c *ClustersGetCommand) Run(parent *Cli, ctx *Context) error {
	resp, err := ctx.APIClient.GetClusterByIdWithResponse(context.TODO(), c.Id)
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Cluster)(resp.JSON200))
}
