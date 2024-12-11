package templates

import (
	"html/template"
	"path/filepath"
)

// Функции для шаблонов
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"minus": func(a, b int) int {
		return b - a
	},
}

func ParseTemplate(fileName string) *template.Template {
	return template.Must(template.ParseFiles(filepath.Join("../../web/templates", fileName)))
}

func ParseTemplateWithFuncs(fileName string) *template.Template {
	return template.Must(template.New("songs.html").Funcs(templateFuncs).ParseFiles("../../web/templates/songs.html"))
}
