package issue

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/urfave/cli"
)

const (
	layout     = "01/02/2006"
	fileLayout = "01-02-2006-150405"
)

func CommentsReport(c *cli.Context) error {
	var fname string
	if len(c.String("output")) > 0 {
		fname = c.String("output")
	} else {
		fname = fmt.Sprintf("%s-%s.csv", c.String("output"), time.Now().Format(fileLayout))
	}
	output, err := os.Create(fname)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to create file %s %s", c.String("output"), err),
			2,
		)
	}
	defer output.Close()
	writer := csv.NewWriter(output)
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	err = writer.Write([]string{
		"Issue ID", "Title", "Total Comments",
		"Status", "Created On", "Closed On",
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	count, err := writeIssues(c, gclient, writer)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	logger.GetLogger(c).Infof("wrote %d records in the report", count)
	return nil
}

func writeIssues(c *cli.Context, gclient *github.Client, writer *csv.Writer) (int, error) {
	opt := &github.IssueListByRepoOptions{
		State:       c.String("state"),
		Sort:        "comments",
		ListOptions: github.ListOptions{PerPage: 30},
	}
	count := 0
	for {
		issues, resp, err := gclient.Issues.ListByRepo(
			context.Background(),
			c.GlobalString("owner"),
			c.GlobalString("repository"),
			opt,
		)
		if err != nil {
			return count, fmt.Errorf("error in fetching issues %s", err)
		}
		for _, iss := range issues {
			if iss.IsPullRequest() {
				continue
			}
			var closedStr string
			if iss.GetState() == "closed" {
				closedStr = iss.GetClosedAt().Format(layout)
			}
			err := writer.Write([]string{
				strconv.Itoa(iss.GetNumber()),
				iss.GetTitle(),
				strconv.Itoa(iss.GetComments()),
				iss.GetState(),
				iss.GetCreatedAt().Format(layout),
				closedStr,
			})
			if err != nil {
				return count, fmt.Errorf("error in writing issues to csv file %s", err)
			}
			count++
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return count, fmt.Errorf("error with csv writing %s", err)
	}
	return count, nil
}
