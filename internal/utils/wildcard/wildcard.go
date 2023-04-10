package wildcard

import (
	"bytes"
	"text/template"
)

func ParseWildCard(text string, wildcard map[string]string) (string, error) {
	parser, err := template.New("wildcard").Parse(text)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = parser.Execute(&tpl, wildcard)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
