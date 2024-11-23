package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type PortalList []flightdeckv1.Portal

func (t PortalList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name"}, table.RowConfig{})
	for _, tenant := range t {
		tw.AppendRow(table.Row{highlight(tenant.Id.String()), tenant.Name})
	}
	return tw
}

type Portal flightdeckv1.Portal

func (t Portal) TableWriter() table.Writer {
	url := t.Url
	if url == "" {
		url = "<pending>"
	}
	data := orderedmap.New()
	data.Set("id", t.Id)
	data.Set("name", t.Name)
	data.Set("title", t.Title)
	data.Set("organization name", t.OrganizationName)
	data.Set("domain", t.Domain)
	data.Set("version", t.Version.Version)
	data.Set("url", url)
	data.Set("arctir tenant", t.TenantName)
	data.Set("alt. domains", t.AlternateDomains)
	return resourceTable(data)
}
