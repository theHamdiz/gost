// handlers/handlers.go
package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ExecuteTemplate(w, "home.templ", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
