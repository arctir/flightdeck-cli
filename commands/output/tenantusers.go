package output

import (
	flightdeckv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

type TenantUserList []flightdeckv1.TenantUser

func (t TenantUserList) TableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{"ID", "Name"})
	for _, user := range t {
		tw.AppendRow(table.Row{highlight(user.Id.String()), user.Username})
	}
	return tw
}

type TenantUser flightdeckv1.TenantUser

func (t TenantUser) TableWriter() table.Writer {
	data := orderedmap.New()
	data.Set("id", t.Id)
	data.Set("username", t.Username)
	data.Set("email", t.Email)
	return resourceTable(data)
}
