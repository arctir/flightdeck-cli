package commands

type CreateCommand struct {
	Tenant           TenantsCreateCommand           `cmd:"tenant" name:"tenant"`
	TenantUser       TenantUsersCreateCommand       `cmd:"tenantuser" name:"tenantuser"`
	IdentityProvider IdentityProvidersCreateCommand `cmd:"identityprovider" name:"identityprovider"`
	Portal           PortalsCreateCommand           `cmd:"portal" name:"portal"`
	Integration      IntegrationsCreateCommand      `cmd:"integration" name:"integration"`
	CatalogProvider  CatalogProvidersCreateCommand  `cmd:"catalogprovider" name:"catalogprovider"`
	Connection       ConnectionsCreateCommand       `cmd:"connection" name:"connection"`
}

type GetCommand struct {
	Clusters          ClustersGetCommand          `cmd:"clusters"`
	Orgs              OrgsGetCommand              `cmd:"orgs"`
	PortalVersions    PortalVersionsGetCommand    `cmd:"portalversions" name:"portalversions"`
	Tenants           TenantsGetCommand           `cmd:"tenants" name:"tenants"`
	TenantUsers       TenantUsersGetCommand       `cmd:"tenantusers" name:"tenantusers"`
	IdentityProviders IdentityProvidersGetCommand `cmd:"identityproviders" name:"identityproviders"`
	Portals           PortalsGetCommand           `cmd:"portals" name:"portals"`
	CatalogProviders  CatalogProvidersGetCommand  `cmd:"catalogproviders" name:"catalogproviders"`
	Integrations      IntegrationsGetCommand      `cmd:"integrations" name:"integrations"`
	Connections       ConnectionsGetCommand       `cmd:"connections" name:"connections"`
}

type DeleteCommand struct {
	Org               OrgsDeleteCommand             `cmd:"org" name:"org"`
	Tenants           TenantsDeleteCommand          `cmd:"tenant" name:"tenant"`
	TenantUsers       TenantUsersDeleteCommand      `cmd:"tenantuser" name:"tenantuser"`
	IdentityProviders IdentityProviderDeleteCommand `cmd:"identityprovider" name:"identityprovider"`
	Portals           PortalDeleteCommand           `cmd:"portals" name:"portals"`
	CatalogProviders  CatalogProvidersDeleteCommand `cmd:"catalogproviders" name:"catalogproviders"`
	Integrations      IntegrationsDeleteCommand     `cmd:"integrations" name:"integrations"`
	Connections       ConnectionsDeleteCommand      `cmd:"connections" name:"connections"`
}
