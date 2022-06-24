package analytics

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	gofn "github.com/repeale/fp-go"
	"github.com/urfave/cli"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

func generateReportRequest(clt *cli.Context) *ga.ReportRequest {
	endDate := time.Now().Format("2006-01-02")
	if len(clt.String("end-date")) > 0 {
		endDate = clt.String("end-date")
	}

	return &ga.ReportRequest{
		ViewId: clt.String("view-id"),
		DateRanges: []*ga.DateRange{
			{EndDate: endDate, StartDate: clt.String("start-date")},
		},
		Dimensions: []*ga.Dimension{
			{Name: "ga:Date"},
		},
		Metrics: []*ga.Metric{
			{Expression: "ga:sessions"},
			{Expression: "ga:pageviews"},
			{Expression: "ga:users"},
		},
	}
}

func Report(clt *cli.Context) error {
	srv, err := ga.NewService(context.Background(), option.WithCredentialsFile(clt.String("credential-file")))
	if err != nil {
		return fmt.Errorf("error in creating service client %s", err)
	}
	rq := &ga.GetReportsRequest{ReportRequests: []*ga.ReportRequest{generateReportRequest(clt)}}
	res, err := ga.NewReportsService(srv).BatchGet(rq).Do()
	if err != nil {
		return fmt.Errorf("error in running the query %s", err)
	}

	return writeOutput(clt, res)
}

func writeOutput(clt *cli.Context, res *ga.GetReportsResponse) error {
	fhr := os.Stdout
	if len(clt.String("output")) > 1 {
		ch, err := os.Create(clt.String("output"))
		if err != nil {
			return fmt.Errorf("error in creating file %s %s", clt.String("output"), err)
		}
		fhr = ch
	}
	defer fhr.Close()
	wrt := csv.NewWriter(fhr)
	err := wrt.Write(processReportHeader(res))
	if err != nil {
		return fmt.Errorf("error in writing header %s", err)
	}
	for _, row := range res.Reports[0].Data.Rows {
		for _, metric := range row.Metrics {
			dataRow := slices.Insert(metric.Values, 0, fmtDate(row.Dimensions[0]))
			if err = wrt.Write(dataRow); err != nil {
				return fmt.Errorf("error in writing data %s", err)
			}
		}
	}
	wrt.Flush()
	if err := wrt.Error(); err != nil {
		return fmt.Errorf("error in finishing csv output %s", err)
	}

	return nil
}

func fmtDate(date string) string {
	parsed, _ := time.Parse("20060102", date)

	return parsed.Format("2006-01-02")
}

func processReportHeader(res *ga.GetReportsResponse) []string {
	hentries := res.Reports[0].ColumnHeader.MetricHeader.MetricHeaderEntries
	dim := res.Reports[0].ColumnHeader.Dimensions[0]
	pipeline := gofn.Pipe3(gofn.Map(metricName), gofn.Map(removeGAprefix), gofn.Map(ucFirst))

	return slices.Insert(pipeline(hentries), 0, removeGAprefix(dim))
}

func ucFirst(s string) string {
	engCase := cases.Title(language.English)

	return engCase.String(s)
}

func metricName(header *ga.MetricHeaderEntry) string {
	return header.Name
}

func removeGAprefix(name string) string {
	return strings.TrimPrefix(name, "ga:")
}
