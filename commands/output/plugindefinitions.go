package output

import (
	"strconv"

	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type PluginDefinitionList []flightdeckv1.PluginDefinition

func (d PluginDefinitionList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name", "Version"}, table.RowConfig{})
	for _, definition := range d {
		tw.AppendRow(table.Row{highlight(definition.Id.String()), definition.Name, definition.Version})
	}
	return tw
}

type PluginDefinition flightdeckv1.PluginDefinition

func (d PluginDefinition) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", d.Id)
	data.Set("name", d.Name)
	data.Set("version", d.Version)
	data.Set("display name", d.DisplayName)
	data.Set("frontend", strconv.FormatBool(d.Frontend))
	data.Set("backend", strconv.FormatBool(d.Backend))
	return resourceTable(data)
}
