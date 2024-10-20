package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type IdentityProviderList []flightdeckv1.IdentityProvider

func (d IdentityProviderList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name"}, table.RowConfig{})
	for _, definition := range d {
		tw.AppendRow(table.Row{highlight(definition.Id.String()), definition.Name})
	}
	return tw
}

type IdentityProvider flightdeckv1.IdentityProvider

func (d IdentityProvider) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", d.Id)
	data.Set("name", d.Name)
	return resourceTable(data)
}