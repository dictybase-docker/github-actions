package analytics

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	ga "google.golang.org/api/analyticsreporting/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestFmtDate(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "2024-01-15", fmtDate("20240115"))
	assert.Equal(t, "0001-01-01", fmtDate("invalid"))
}

func TestUcFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"sessions", "Sessions"},
		{"pageviews", "Pageviews"},
		{"users", "Users"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ucFirst(tt.input))
		})
	}
}

func TestRemoveGAprefix(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "date", removeGAprefix("ga:date"))
	assert.Equal(t, "sessions", removeGAprefix(sessionsMetric))
	assert.Equal(t, "noprefix", removeGAprefix("noprefix"))
}

func TestMetricName(t *testing.T) {
	t.Parallel()

	entry := &ga.MetricHeaderEntry{Name: sessionsMetric}
	assert.Equal(t, sessionsMetric, metricName(entry))
}

func TestProcessReportHeader(t *testing.T) {
	t.Parallel()

	res := &ga.GetReportsResponse{
		Reports: []*ga.Report{{
			ColumnHeader: &ga.ColumnHeader{
				Dimensions: []string{"ga:date"},
				MetricHeader: &ga.MetricHeader{
					MetricHeaderEntries: []*ga.MetricHeaderEntry{
						{Name: sessionsMetric},
						{Name: "ga:pageviews"},
						{Name: "ga:users"},
					},
				},
			},
		}},
	}
	result := processReportHeader(res)
	assert.Equal(t, []string{"date", "Sessions", "Pageviews", "Users"}, result)
}

func TestGenerateReportRequest(t *testing.T) {
	t.Parallel()

	t.Run("with default dates", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("view-id", "12345", "")
		flagSet.String("start-date", "2024-01-01", "")
		flagSet.String("end-date", "", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "view-id"},
			cli.StringFlag{Name: "start-date"},
			cli.StringFlag{Name: "end-date"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		req := generateReportRequest(ctx)
		assert.Equal(t, "12345", req.ViewId)
		assert.Len(t, req.DateRanges, 1)
		assert.Equal(t, "2024-01-01", req.DateRanges[0].StartDate)
		assert.NotEmpty(t, req.DateRanges[0].EndDate)
	})

	t.Run("with custom end date", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("view-id", "12345", "")
		flagSet.String("start-date", "2024-01-01", "")
		flagSet.String("end-date", "2024-06-30", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "view-id"},
			cli.StringFlag{Name: "start-date"},
			cli.StringFlag{Name: "end-date"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		req := generateReportRequest(ctx)
		assert.Equal(t, "2024-06-30", req.DateRanges[0].EndDate)
	})
}

func TestWriteOutput(t *testing.T) {
	t.Parallel()

	res := &ga.GetReportsResponse{
		Reports: []*ga.Report{{
			ColumnHeader: &ga.ColumnHeader{
				Dimensions: []string{"ga:date"},
				MetricHeader: &ga.MetricHeader{
					MetricHeaderEntries: []*ga.MetricHeaderEntry{
						{Name: sessionsMetric},
					},
				},
			},
			Data: &ga.ReportData{
				Rows: []*ga.ReportRow{
					{
						Dimensions: []string{"20240101"},
						Metrics:    []*ga.DateRangeValues{{Values: []string{"100"}}},
					},
					{
						Dimensions: []string{"20240102"},
						Metrics:    []*ga.DateRangeValues{{Values: []string{"200"}}},
					},
				},
			},
		}},
	}

	t.Run("write to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outFile := filepath.Join(tmpDir, "output.csv")

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("output", outFile, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{cli.StringFlag{Name: "output"}}

		ctx := cli.NewContext(app, flagSet, nil)
		err := writeOutput(ctx, res)
		require.NoError(t, err)

		data, err := os.ReadFile(outFile)
		require.NoError(t, err)
		assert.Contains(t, string(data), "date,Sessions")
		assert.Contains(t, string(data), "2024-01-01,100")
		assert.Contains(t, string(data), "2024-01-02,200")
	})

	t.Run("write to tmp file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outFile := filepath.Join(tmpDir, "report.csv")

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("output", outFile, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{cli.StringFlag{Name: "output"}}

		ctx := cli.NewContext(app, flagSet, nil)
		err := writeOutput(ctx, res)
		require.NoError(t, err)
		assert.FileExists(t, outFile)
	})
}
