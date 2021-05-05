package comment

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/v32/github"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/ontology"

	"github.com/urfave/cli"
)

const tmpl = `
{{if index . "fail"}}
## :x: :x: Ontology errors :x: :x:
{{- range $i,$rc := index . "fail"}}
> ### File: *{{$rc.Name}}*
:boom:**Violations**:boom:
{{- range $y,$v := $rc.Violations}}
- {{$v}}
{{- end}}
{{- end}}
{{- end}}

{{ if index . "pass" -}}
----

## :heavy_check_mark: :heavy_check_mark: Ontology pass :heavy_check_mark: :heavy_check_mark:
{{- range $i,$rc := index . "pass"}}
> ### File: *{{$rc.Name}}*
{{- end}}
{{- end}}
`

type reportParams struct {
	data       map[string][]*reportContent
	owner      string
	repository string
	token      string
	prid       int
	ref        string
}

type checkStatusParams struct {
	data       map[string][]*reportContent
	client     *github.Client
	owner      string
	repository string
	ref        string
	report     string
}

type reportContent struct {
	Name       string
	Violations []string
}

func OntoReportOnPullComment(c *cli.Context) error {
	cf, err := listCommittedFiles(c.String("commit-list-file"))
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	rs := make(map[string][]*reportContent)
	for _, f := range cf {
		v, err := ontology.ParseViolations(
			fmt.Sprintf("%s/%s.json", c.String("report-dir"), f),
			"ERROR",
		)
		if err != nil {
			if !ontology.IsViolationNotFound(err) {
				return cli.NewExitError(err.Error(), 2)
			}
			if _, ok := rs["pass"]; ok {
				rs["pass"] = append(rs["pass"],
					&reportContent{Name: fmt.Sprintf("%s.obo", f)},
				)
			} else {
				rs["pass"] = []*reportContent{{Name: fmt.Sprintf("%s.obo", f)}}
			}
			continue
		}
		if _, ok := rs["fail"]; ok {
			rs["fail"] = append(rs["fail"],
				&reportContent{Name: fmt.Sprintf("%s.obo", f), Violations: v},
			)
			continue
		}
		rs["fail"] = []*reportContent{{Name: fmt.Sprintf("%s.obo", f), Violations: v}}
	}
	err = createCommentFromReport(&reportParams{
		prid:       c.Int("pull-request-id"),
		repository: c.GlobalString("repository"),
		owner:      c.GlobalString("owner"),
		token:      c.GlobalString("token"),
		ref:        c.String("ref"),
		data:       rs,
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return reportStatusError(rs)
}

func reportStatusError(rs map[string][]*reportContent) error {
	if _, ok := rs["fail"]; ok {
		return cli.NewExitError(
			fmt.Sprintf("failed report count %d", len(rs["fail"])),
			2,
		)
	}
	return nil
}

func createCommentFromReport(args *reportParams) error {
	gclient, err := client.GetGithubClient(args.token)
	if err != nil {
		return err
	}
	mk, err := mkdownOutput(args.data)
	if err != nil {
		return err
	}
	_, _, err = gclient.Issues.CreateComment(
		context.Background(),
		args.owner,
		args.repository,
		args.prid,
		&github.IssueComment{
			Body: github.String(mk.String()),
		})
	if err != nil {
		return fmt.Errorf("error in creating pull request comment %s", err)
	}
	return nil
}

func mkdownOutput(data interface{}) (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	t, err := template.New("onto-report").Parse(tmpl)
	if err != nil {
		return out, fmt.Errorf("error in parsing template %s", err)
	}
	if err := t.Execute(out, data); err != nil {
		return out, fmt.Errorf("error in executing template %s", err)
	}
	return out, nil
}

func listCommittedFiles(path string) ([]string, error) {
	var a []string
	r, err := os.Open(path)
	if err != nil {
		return a, fmt.Errorf("unable to open file %s", err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		a = append(a, baseNoSuffix(scanner.Text()))
	}
	return a, scanner.Err()
}

func baseNoSuffix(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}

func manageCheckStatus(args *checkStatusParams) error {
	chresult, _, err := args.client.Checks.ListCheckRunsForRef(
		context.Background(),
		args.owner,
		args.repository,
		args.ref,
		&github.ListCheckRunsOptions{
			CheckName: github.String(args.report),
		},
	)
	if err != nil {
		return fmt.Errorf("error in listing check runs %s", err)
	}
	if chresult.GetTotal() < 1 {
		return fmt.Errorf("should have at least one check status, got %d", chresult.GetTotal())
	}
	return updateCheckStatus(chresult.CheckRuns[0], args)
}

func updateCheckStatus(res *github.CheckRun, args *checkStatusParams) error {
	concl := "success"
	if _, ok := args.data["fail"]; ok {
		concl = "failure"
	}
	_, _, err := args.client.Checks.UpdateCheckRun(
		context.Background(),
		args.owner,
		args.repository,
		res.GetID(),
		github.UpdateCheckRunOptions{
			Name:       args.report,
			Status:     github.String("completed"),
			Conclusion: github.String(concl),
		},
	)
	if err != nil {
		return fmt.Errorf("error in updating status check %s", err)
	}
	return nil
}
