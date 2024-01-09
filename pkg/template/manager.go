package template

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

var ErrNoTemplate = errors.New("no template")
var ErrExecTemplate = errors.New("exec template")

type Manager struct {
	// filename ~ template name
	templates map[string]*template.Template
}

func NewManager() *Manager {
	return &Manager{
		templates: make(map[string]*template.Template),
	}
}

func (m *Manager) ParseTemplates(templates map[string][]string) error {
	for name, filenames := range templates {
		for i := range filenames {
			filenames[i] = fmt.Sprintf("templates/%s.html", filenames[i])
		}

		templ, err := template.ParseFiles(
			filenames...,
		)
		if err != nil {
			return err
		}

		if _, ok := m.templates[name]; !ok {
			m.templates[name] = templ
		}
	}

	return nil
}

func (m *Manager) Execute(w http.ResponseWriter, name string, obj interface{}) {
	templ, ok := m.templates[name]
	if !ok {
		writeError(w, ErrNoTemplate)
		return
	}

	if err := templ.Execute(w, obj); err != nil {
		writeError(w, ErrExecTemplate)
	}
}
