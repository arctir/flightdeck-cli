package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type TenantList []flightdeckv1.Tenant

func (t TenantList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name"})
	for _, tenant := range t {
		tw.AppendRow(table.Row{highlight(tenant.Id.String()), tenant.Name})
	}
	return tw
}

type Tenant flightdeckv1.Tenant

func (t Tenant) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", t.Id)
	data.Set("name", t.Name)
	data.Set("display name", t.DisplayName)
	data.Set("identifier", t.Identifier)
	return resourceTable(data)
}
