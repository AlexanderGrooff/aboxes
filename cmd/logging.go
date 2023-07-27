package cmd

import (
	"bytes"
	"text/template"
)

// Type for results
type Result struct {
	Target   string
	Hostname string
	Stdout   string
	Stderr   string
	Error    error
}

func (result *Result) toString(format string) string {
	tmpl, err := template.New("result").Parse(format)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, result)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
