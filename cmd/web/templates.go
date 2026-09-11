package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.alexedwards.net/internal/models"
)

// FuncMap functions must return only ONE thing. You can have err as second, but nothing else.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl.html". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl.html ui/html/pages/view.tmpl.html]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl.html') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil { 
			return nil, err
		}

		// Call ParseFilse() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.htmpl.html') as the key.
		cache[name] = ts
	}

	return cache, nil
}

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
type templateData struct {
	CurrentYear int
	Snippet models.Snippet
	Snippets []models.Snippet
}