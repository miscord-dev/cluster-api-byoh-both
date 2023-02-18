package machinetemplate

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

func NewTemplate(tmpl string) (Template, error) {
	findStringSubmatch := func(expr, value string) []string {
		reg := regexp.MustCompile(expr)

		return reg.FindStringSubmatch(value)
	}

	template, err := template.New("").Funcs(template.FuncMap{
		"join":               strings.Join,
		"split":              strings.Split,
		"contains":           strings.Contains,
		"hasPrefix":          strings.HasPrefix,
		"hasSuffix":          strings.HasSuffix,
		"toLower":            strings.ToLower,
		"toUpper":            strings.ToUpper,
		"findStringSubmatch": findStringSubmatch,
	}).Parse(tmpl)

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &textTemplate{
		tmpl: template,
	}, nil
}

type NodeInfo struct {
	OS         string
	OSImage    string
	Arch       string
	K8sVersion string
}

type ScriptInfo struct {
	ImageTag           string
	BundleDownloadPath string
	NodeInfo
}

type Template interface {
	GenerateTag(nodeInfo NodeInfo) (string, error)
	GenerateScript(scriptInfo ScriptInfo) (string, error)
}

type textTemplate struct {
	tmpl *template.Template
}

func (t *textTemplate) generate(value any) (_ string, reterr error) {
	defer func() {
		err := recover()

		if err != nil {
			reterr = err.(error)
		}
	}()

	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, value)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *textTemplate) GenerateTag(nodeInfo NodeInfo) (_ string, reterr error) {
	value, err := t.generate(nodeInfo)

	if err != nil {
		return "", fmt.Errorf("failed to generate tag: %w", err)
	}

	return value, nil
}

func (t *textTemplate) GenerateScript(scriptInfo ScriptInfo) (string, error) {
	value, err := t.generate(scriptInfo)

	if err != nil {
		return "", fmt.Errorf("failed to generate script: %w", err)
	}

	return value, nil
}
