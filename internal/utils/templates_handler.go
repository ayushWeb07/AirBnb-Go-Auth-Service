package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

type EmailData struct {
	UserName string
	AppName  string
	Otp      string
}

func Render(name string, data any) (string, error) {
	cwd, _ := os.Getwd()
	var templates = template.Must(template.ParseGlob(filepath.Join(cwd, "internal/utils/templates/mail/*.txt")))

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}
