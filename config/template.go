package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type SheTemplate struct {
	tmpl *template.Template
}

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func funcMap() template.FuncMap {
	result := make(template.FuncMap)
	result["prompt"] = func(label string) string {
		return stringPrompt(label)
	}

	return result
}

func NewTemplate() *SheTemplate {
	result := SheTemplate{
		tmpl: template.New("SheTemplate").Funcs(sprig.FuncMap()).Funcs(funcMap()),
	}

	return &result
}

func (t *SheTemplate) Execute(s string, data any) (string, error) {
	tmpl, err := t.tmpl.Clone()
	if err != nil {
		return "", err
	}

	tmpl, err = tmpl.Parse(s)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
