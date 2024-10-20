package server

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func ParseTemplateResponse(name string, input interface{}) (string, error) {
	file, err := templates.ReadFile(fmt.Sprintf("templates/%s.tmpl", name))
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	tmpl := template.Must(template.New(name).Parse(string(file)))
	if err := tmpl.Execute(buf, input); err != nil {
		return "", err
	}
	return buf.String(), nil
}
