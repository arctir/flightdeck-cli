package commands

import (
	apiv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	flightdeckclient "github.com/arctir/go-flightdeck/pkg/client"
)

type Context struct {
	APIClient       *apiv1.ClientWithResponses
	Config          *flightdeckclient.Config
	SkipConfigCheck bool
}
