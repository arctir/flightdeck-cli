package output

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	orderedmap "github.com/wk8/go-ordered-map"
)

func IsOutputtableResponse(response *http.Response) bool {
	switch response.StatusCode {
	case http.StatusNotFound,
		http.StatusForbidden,
		http.StatusBadRequest,
		http.StatusConflict,
		http.StatusInternalServerError,
		http.StatusUnauthorized:
		fmt.Println(http.StatusText(response.StatusCode))
		return false
	case http.StatusNoContent:
		fmt.Println("Success")
		return false
	}
	return true
}

func OutputResult[T apiResponse](format string, t T) error {
	var err error
	if format == "json" {
		err = outputJson(t)
	} else {
		err = outputTable(t)
	}
	return err
}

type apiResponse interface {
	*Cluster | *ClusterList |
		*Organization | *OrganizationList |
		*Portal | *PortalList |
		*CatalogProvider | *CatalogProviderList |
		*Integration | *IntegrationList |
		*Tenant | *TenantList |
		*TenantUser | *TenantUserList |
		*IdentityProvider | *IdentityProviderList |
		*PortalVersion | *PortalVersionList |
		*Connection | *ConnectionList

	TableWriter() table.Writer
}

func outputJson[T apiResponse](t T) error {
	s, err := json.MarshalIndent(&t, "", "   ")
	if err != nil {
		return err
	}
	fmt.Println(string(s))
	return nil
}

func outputTable[T apiResponse](t T) error {
	tw := t.TableWriter()
	fmt.Println(tw.Render())
	return nil
}

func highlight(t string) string {
	return fmt.Sprintf("%s%s%s%s", text.Bold.EscapeSeq(), text.FgHiGreen.EscapeSeq(), t, text.Reset.EscapeSeq())
}

func resourceTable(data *orderedmap.OrderedMap) table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower
	tw.AppendHeader(table.Row{text.AlignRight.Apply("attribute", 16), "value"})
	for pair := data.Oldest(); pair != nil; pair = pair.Next() {
		tw.AppendRow(table.Row{text.AlignRight.Apply(highlight(pair.Key.(string)), 16), pair.Value})
	}
	return tw
}
