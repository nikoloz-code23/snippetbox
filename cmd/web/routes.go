package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes(cfg *config) *http.ServeMux {
	mux := http.NewServeMux()

	// Point direction to the file server location. Returns http.Handler.
	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	// Use the mux.Handle() function to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we strip
	// the "/static" prefix before the request reaches the file server.

	//http.StripPrefix Returns http.Handler.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", 														 app.home)
	mux.HandleFunc("GET /snippet/view/{id}", 		app.snippetView)
	mux.HandleFunc("GET /snippet/create", 				app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", 		app.snippetCreatePost)

	return mux
}