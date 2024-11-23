package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type PortalVersionList []flightdeckv1.PortalVersion

func (l PortalVersionList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Version"})
	for _, org := range l {
		tw.AppendRow(table.Row{highlight(org.Id.String()), org.Version})
	}
	return tw
}

type PortalVersion flightdeckv1.PortalVersion

func (t PortalVersion) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", t.Id)
	data.Set("version", t.Version)
	data.Set("major", t.Major)
	data.Set("minor", t.Minor)
	data.Set("patch", t.Patch)
	data.Set("revision", t.Rev)
	return resourceTable(data)
}
