package common

import "github.com/google/uuid"

type OrgFlags struct {
	Org uuid.UUID `name:"org" short:"o" default:"${defaultOrg}" help:"ID of the organization."`
}

type PortalFlags struct {
	PortalName string `arg:"portal-name" help:"Name of the portal."`
	OrgFlags
}

type TenantFlags struct {
	TenantName string `arg:"tenant-name" help:"Name of the tenant."`
	OrgFlags
}
