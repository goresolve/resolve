package resolve

import (
	"html/template"
)

type Template struct {
	file     string
	template *template.Template
}

func NewTemplate(name, file string) Template {
	tmpl, err := template.New(name).ParseFiles(file)
	if err != nil {
		ErrorMessage(err.Error(), 3)
	}

	return Template{
		file:     file,
		template: tmpl,
	}
}

func (t *Template) Reload() {
	tmpl, err := template.New(t.template.Name()).ParseFiles(t.file)
	if err != nil {
		panic(err)
	}

	t.template = tmpl
}
