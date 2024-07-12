package issue

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v62/github"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/urfave/cli"
)

const (
	layout           = "01/02/2006"
	dateFilterlayout = "2006-01-02"
	fileLayout       = "01-02-2006-150405"
)

func CommentsCountByDate(clt *cli.Context) error {
	parsedDate, err := time.Parse(dateFilterlayout, clt.String("since"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in parsing since parameter %s", err),
			2,
		)
	}
	gclient := github.NewClient(nil)
	opt := &github.IssueListByRepoOptions{
		Since:       parsedDate,
		ListOptions: github.ListOptions{PerPage: 15},
	}
	var totalComments, totalIssues int
	for {
		issues, resp, err := gclient.Issues.ListByRepo(
			context.Background(),
			clt.GlobalString("owner"),
			clt.GlobalString("repository"),
			opt,
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error in fetching issues %s", err),
				2,
			)
		}
		totalIssues += len(issues)
		for _, iss := range issues {
			totalComments += *iss.Comments
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Printf("total no of issues %d\n", totalIssues)
	fmt.Printf("total no of comments %d\n", totalComments)
	return nil
}

func CommentsReport(clt *cli.Context) error {
	var fname string
	if len(clt.String("output")) > 0 {
		fname = clt.String("output")
	} else {
		fname = fmt.Sprintf("%s-%s.csv", clt.String("output"), time.Now().Format(fileLayout))
	}
	output, err := os.Create(fname)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf(
				"unable to create file %s %s",
				clt.String("output"),
				err,
			),
			2,
		)
	}
	defer output.Close()
	writer := csv.NewWriter(output)
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	count, err := writeIssues(clt, gclient, writer)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	logger.GetLogger(clt).Infof("wrote %d records in the report", count)

	return nil
}

func writeIssues(
	clt *cli.Context,
	gclient *github.Client,
	writer *csv.Writer,
) (int, error) {
	count := 0
	err := writer.Write([]string{
		"Issue ID", "Title", "Total Comments",
		"Status", "Created On", "Closed On",
	})
	if err != nil {
		return count, fmt.Errorf("error in writing file header %s", err)
	}
	opt := issueOpts(clt)
	for {
		issues, resp, err := gclient.Issues.ListByRepo(
			context.Background(),
			clt.GlobalString("owner"),
			clt.GlobalString("repository"),
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
				return count, fmt.Errorf(
					"error in writing issues to csv file %s",
					err,
				)
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
		return count, fmt.Errorf("error in writing %s", err)
	}

	return count, nil
}

func issueOpts(c *cli.Context) *github.IssueListByRepoOptions {
	return &github.IssueListByRepoOptions{
		State:       c.String("state"),
		Sort:        "comments",
		ListOptions: github.ListOptions{PerPage: 30},
	}
}
