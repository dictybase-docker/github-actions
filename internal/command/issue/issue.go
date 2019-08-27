package issue

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/google/go-github/v28/github"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/urfave/cli"
)

const (
	layout = "01/02/2006"
)

func CommentsReport(c *cli.Context) error {
	var output io.Writer
	if len(c.String("output")) > 0 {
		wf, err := os.Create(c.String("output"))
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("unable to open file %s %s", c.String("output"), err),
				2,
			)
		}
		output = wf
	} else {
		output = os.Stdout
	}
	writer := csv.NewWriter(output)
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	opt := &github.IssueListByRepoOptions{
		State:       c.String("state"),
		Sort:        "comments",
		ListOptions: github.ListOptions{PerPage: 30},
	}
	writer.Write([]string{
		"Issue ID", "Title", "Total Comments",
		"Status", "Created On", "Closed On",
	})
	for {
		issues, resp, err := gclient.Issues.ListByRepo(
			context.Background(),
			c.GlobalString("owner"),
			c.GlobalString("repository"),
			opt,
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error in fetching issues %s", err),
				2,
			)
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
				return cli.NewExitError(
					fmt.Sprintf("error in writing issues to csv file %s", err),
					2,
				)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error with csv writing %s", err),
			2,
		)
	}
	return nil
}
