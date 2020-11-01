package views

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var (
	LayoutDir   = "views/layout/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func layoutFiles() []string {
	wd, _ := os.Getwd()
	files, err := filepath.Glob(wd + "/pkg/" + LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// A helper function that panics on any error
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
