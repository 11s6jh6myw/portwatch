package trend

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Report writes trend summaries to w in the given format ("text" or "json").
func Report(w io.Writer, trends []PortTrend, format string) error {
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Port < trends[j].Port
	})
	switch format {
	case "json":
		return json.NewEncoder(w).Encode(trends)
	default:
		return reportText(w, trends)
	}
}

func reportText(w io.Writer, trends []PortTrend) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tOPENED\tCLOSED\tFLAPPING\tLAST SEEN")
	for _, t := range trends {
		flap := "no"
		if t.Flapping {
			flap = "YES"
		}
		fmt.Fprintf(tw, "%d\t%d\t%d\t%s\t%s\n",
			t.Port, t.OpenCount, t.CloseCount, flap,
			t.LastSeen.Format("2006-01-02 15:04:05"))
	}
	return tw.Flush()
}
