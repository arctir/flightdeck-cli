package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type OrganizationList []flightdeckv1.Organization

func (l OrganizationList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name"})
	for _, org := range l {
		tw.AppendRow(table.Row{highlight(org.Id.String()), org.Name})
	}
	return tw
}

type Organization flightdeckv1.Organization

func (t Organization) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", t.Id)
	data.Set("name", t.Name)
	data.Set("subdomain", t.Subdomain)
	return resourceTable(data)
}
