package common

import "github.com/google/uuid"

type OrgFlags struct {
	Org uuid.UUID `arg:"org" help:"ID of the organization."`
}

type PortalFlags struct {
	OrgFlags
	PortalName string `arg:"portal-name" help:"Name of the portal."`
}

type TenantFlags struct {
	OrgFlags
	TenantName string `arg:"tenant-name" help:"Name of the tenant."`
}
