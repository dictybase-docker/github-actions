package issue

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/email"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	parser "github.com/dictyBase-docker/github-actions/internal/parser"
	"github.com/google/go-github/v62/github"
	"github.com/urfave/cli"
	"golang.org/x/net/html"
)

const (
	layout           = "01/02/2006"
	dateFilterlayout = "2006-01-02"
	fileLayout       = "01-02-2006-150405"
)

func CommentsCountByDate(clt *cli.Context) error {
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 15},
	}
	var totalComments, totalIssues int
	query := fmt.Sprintf(
		"repo:%s/%s created:>=%s",
		clt.GlobalString("owner"),
		clt.GlobalString("repository"),
		clt.String("since"),
	)
	for {
		result, resp, err := gclient.Search.Issues(
			context.Background(),
			query,
			opt,
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error in fetching issues %s", err),
				2,
			)
		}
		totalIssues += len(result.Issues)
		for _, iss := range result.Issues {
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

func getIssue(gclient *github.Client, clt *cli.Context) (*github.Issue, error) {
	// Get issue number from context
	issueNumber := clt.Int("issueid")
	if issueNumber == 0 {
		return nil, fmt.Errorf("issue number is required")
	}

	// Get issue using the Issues API
	issue, _, err := gclient.Issues.Get(
		context.Background(),
		clt.GlobalString("owner"),
		clt.GlobalString("repository"),
		issueNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching issue: %w", err)
	}

	return issue, nil
}

func getIssueBody(issue *github.Issue) (string, error) {
	body := issue.GetBody()
	if body == "" {
		return "", fmt.Errorf("issue body is empty")
	}

	return body, nil
}

func fetchAndParseIssue(clt *cli.Context) (*html.Node, error) {
	// Get GitHub client
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return nil, cli.NewExitError(
			fmt.Sprintf("error getting github client: %s", err),
			2,
		)
	}

	// Fetch the issue
	issue, err := getIssue(gclient, clt)
	if err != nil {
		return nil, cli.NewExitError(
			fmt.Sprintf("error fetching issue: %s", err),
			2,
		)
	}

	// Get the markdown body
	markdownBody, err := getIssueBody(issue)
	if err != nil {
		return nil, cli.NewExitError(
			fmt.Sprintf("error getting issue body: %s", err),
			2,
		)
	}

	// Convert markdown to HTML node
	htmlNode, err := parser.MarkdownToHTML(markdownBody)
	if err != nil {
		return nil, cli.NewExitError(
			fmt.Sprintf("error converting markdown to HTML: %s", err),
			2,
		)
	}

	return htmlNode, nil
}

func SendIssueLabelEmail(clt *cli.Context) error {
	htmlNode, err := fetchAndParseIssue(clt)
	if err != nil {
		return err
	}

	// Extract order data using parser
	issueData, err := parser.ExtractOrderData(htmlNode)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error extracting order data: %s", err),
			2,
		)
	}

	// Validate extracted data
	if issueData.RecipientEmail == "" {
		return cli.NewExitError("no recipient email found in issue", 2)
	}
	if issueData.OrderID == "" {
		return cli.NewExitError("no order ID found in issue", 2)
	}

	// Convert to email data structure
	emailData := email.OrderEmailData{
		RecipientEmail: issueData.RecipientEmail,
		OrderID:        issueData.OrderID,
		Label:          clt.String("label"),
	}

	log := logger.GetLogger(clt)
	log.WithFields(map[string]any{
		"order_id":  emailData.OrderID,
		"recipient": emailData.RecipientEmail,
		"label":     emailData.Label,
	}).Info("Extracted order data from issue")

	// Get Mailgun configuration from flags
	domain := clt.String("domain")
	apiKey := clt.String("apiKey")
	fromEmail := clt.String("fromEmail")

	if domain == "" || apiKey == "" {
		return cli.NewExitError(
			"Mailgun domain and apiKey are required",
			2,
		)
	}

	// Create email client and send email
	emailClient := email.NewEmailClient(domain, apiKey, fromEmail)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := emailClient.SendOrderUpdateFromTemplate(ctx, emailData.RecipientEmail, emailData); err != nil {
		log.WithFields(map[string]any{
			"order_id":  emailData.OrderID,
			"recipient": emailData.RecipientEmail,
			"label":     emailData.Label,
			"error":     err.Error(),
		}).Error("Failed to send email")
		return cli.NewExitError(
			fmt.Sprintf("error sending email: %s", err),
			2,
		)
	}

	log.WithFields(map[string]any{
		"order_id":  emailData.OrderID,
		"recipient": emailData.RecipientEmail,
		"label":     emailData.Label,
		"status":    "sent",
	}).Info("Successfully sent order update email")

	return nil
}
