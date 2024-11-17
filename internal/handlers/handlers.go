package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ShowUploadForm(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	tmpl, err := template.ParseFiles(
		"templates/index.html",
		"templates/components/nav.html",
		"templates/components/tabs.html",
		"templates/components/upload.html",
		"templates/components/quick-actions.html",
		"templates/components/options-panel.html",
		"templates/components/progress-result.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Failed to load templates", http.StatusInternalServerError)
		return
	}

	// Render the index template
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
