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
{{if $rc.HTML}}
> #### Full report
{{$rc.HTML}}
{{- end}}
{{- end}}
{{- end}}

{{ if index . "pass" -}}
----

## :heavy_check_mark: :heavy_check_mark: Ontology pass :heavy_check_mark: :heavy_check_mark:
{{- range $i,$rc := index . "pass"}}
> ### File: *{{$rc.Name}}*
{{if $rc.HTML}}
> #### Full report
{{$rc.HTML}}
{{- end}}
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

type reportContent struct {
	Name       string
	HTML       string
	Violations []string
}

func OntoReportOnPullComment(clt *cli.Context) error {
	cf, err := listCommittedFiles(clt.String("commit-list-file"))
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	rps, err := ontoReport(clt, cf)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	err = createCommentFromReport(&reportParams{
		prid:       clt.Int("pull-request-id"),
		repository: clt.GlobalString("repository"),
		owner:      clt.GlobalString("owner"),
		token:      clt.GlobalString("token"),
		ref:        clt.String("ref"),
		data:       rps,
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	return reportStatusError(rps)
}

func ontoReport(
	clt *cli.Context,
	cf []string,
) (map[string][]*reportContent, error) {
	rcs := make(map[string][]*reportContent)
	for _, folder := range cf {
		html, err := readHTMLContent(
			fmt.Sprintf(
				"%s.html",
				filepath.Join(clt.String("report-dir"), folder),
			),
		)
		if err != nil {
			return rcs, err
		}
		viol, err := ontology.ParseViolations(
			fmt.Sprintf("%s/%s.json", clt.String("report-dir"), folder),
			"ERROR",
		)
		if err != nil {
			if !ontology.IsViolationNotFound(err) {
				return rcs, fmt.Errorf("ontology not found %s", err)
			}
			if _, ok := rcs["pass"]; ok {
				rcs["pass"] = append(rcs["pass"], &reportContent{
					Name: fmt.Sprintf("%s.obo", folder),
					HTML: html,
				})
			} else {
				rcs["pass"] = []*reportContent{{Name: fmt.Sprintf("%s.obo", folder), HTML: html}}
			}

			continue
		}
		if _, ok := rcs["fail"]; ok {
			rcs["fail"] = append(rcs["fail"], &reportContent{
				Name:       fmt.Sprintf("%s.obo", folder),
				Violations: viol,
				HTML:       html,
			})

			continue
		}
		rcs["fail"] = []*reportContent{{
			Name:       fmt.Sprintf("%s.obo", folder),
			Violations: viol,
			HTML:       html,
		}}
	}

	return rcs, nil
}

func readHTMLContent(file string) (string, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", nil
	}
	ct, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error in reading file %s", err)
	}

	return string(ct), nil
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
		return fmt.Errorf("error in getting github client %s", err)
	}
	mkd, err := mkdownOutput(args.data)
	if err != nil {
		return err
	}
	_, _, err = gclient.Issues.CreateComment(
		context.Background(),
		args.owner,
		args.repository,
		args.prid,
		&github.IssueComment{
			Body: github.String(mkd.String()),
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
	var afiles []string
	r, err := os.Open(path)
	if err != nil {
		return afiles, fmt.Errorf("unable to open file %s", err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		afiles = append(afiles, baseNoSuffix(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return afiles, fmt.Errorf("error from scanning %s", err)
	}

	return afiles, nil
}

func baseNoSuffix(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}
