package commands

import (
	"context"
	"errors"

	"github.com/arctir/flightdeck-cli/commands/output"
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/google/uuid"
)

type ClustersGetCommand struct {
	Id *uuid.UUID `arg:"id" optional:"" name:"id" help:"ID of the cluster to get. If not provided, lists all clusters."`
}

func (c ClustersGetCommand) Run(parent *Cli, ctx *Context) error {
	if c.Id == nil {
		return c.list(parent, ctx)
	}
	return c.get(parent, ctx)
}

func (c ClustersGetCommand) list(parent *Cli, ctx *Context) error {
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

func (c ClustersGetCommand) get(parent *Cli, ctx *Context) error {
	if c.Id == nil {
		return errors.New("id is required")
	}
	resp, err := ctx.APIClient.GetClusterByIdWithResponse(context.TODO(), c.Id.String())
	if err != nil {
		return err
	}
	if !output.IsOutputtableResponse(resp.HTTPResponse) {
		return nil
	}
	return output.OutputResult(parent.OutputFormat, (*output.Cluster)(resp.JSON200))
}
