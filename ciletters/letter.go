//go:build !solution
// +build !solution

package ciletters

import (
	"strings"
	"text/template"
)

const tpl = `Your pipeline #{{ .Pipeline.ID}} {{if ne .Pipeline.Status "ok"}}has failed{{else}}passed{{end}}!
    Project:      {{ .Project.GroupID }}/{{.Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ slice .Commit.Hash 0 8 }} {{ .Commit.Message }}
    CommitAuthor: {{ .Commit.Author }}{{range $job := .Pipeline.FailedJobs}}
        Stage: {{$job.Stage}}, Job {{$job.Name}}{{range cmdLog $job.RunnerLog}}
            {{.}}{{end}}
{{end}}`

func cmdLog(s string) []string {
	res := strings.Split(s, "\n")
	if len(res) > 9 {
		res = res[9:]
	}
	return res
}
func MakeLetter(n *Notification) (string, error) {

	funcMap := template.FuncMap{
		"cmdLog": cmdLog,
	}

	t, err := template.New("email").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	if err = t.Execute(&sb, n); err != nil {
		return "", err
	}

	return sb.String(), nil
}
