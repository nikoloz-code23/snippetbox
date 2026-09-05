package main

import "snippetbox.alexedwards.net/internal/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
type templateData struct {
	Snippet models.Snippet
}

func newTemplateData(Snippet models.Snippet) (*templateData) {
	var snippetData models.Snippet = Snippet
	snippetData.Created = snippetData.Created.UTC()
	snippetData.Expires = snippetData.Expires.UTC()

	return &templateData{
		Snippet: snippetData,
	}
}