package common

type OrgFlags struct {
	Org string `arg:"org"`
}

type PortalFlags struct {
	OrgFlags
	PortalName string `arg:"portal-name"`
}

type TenantFlags struct {
	OrgFlags
	TenantName string `arg:"portal-name"`
}
