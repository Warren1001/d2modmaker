package server

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/alecthomas/template"
	"github.com/tlentz/d2modmaker/gui/webpack"
)

// Handler returns http.Handler for server endpoint
func Handler(isProduction bool, buildPath string) http.HandlerFunc {
	tmpl, err := parseTemplate(isProduction)

	if err != nil {
		return func(res http.ResponseWriter, req *http.Request) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}

	data, err := NewViewData(isProduction, buildPath)

	if err != nil {
		return func(res http.ResponseWriter, req *http.Request) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if err := tmpl.Execute(res, data); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func parseTemplate(isProduction bool) (*template.Template, error) {
	if isProduction {
		filename := path.Join("templates", "index.html")
		fmt.Println(filename)
		b, err := webpack.ReadFile(isProduction, filename)
		if err != nil {
			log.Fatal(err)
		}
		s := string(b)
		name := filepath.Base(filename)
		t := template.New(name)
		_, err = t.Parse(s)
		return t, err
	}
	return template.ParseFiles(path.Join("templates", "index.html"))
}
