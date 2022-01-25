package main

import (
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/charlesharries/wordle-go/pkg/forms"
)

type templateData struct {
	Attempt         string
	CSRFToken       string
	CurrentYear     int
	DevMode         bool
	Flash           string
	Form            *forms.Form
	FormData        url.Values
	FormErrors      map[string]string
	IsAuthenticated bool
	Game            game
}

func times(times int) []interface{} {
	return make([]interface{}, times)
}

func sub(first, second int) int {
	return first - second
}

// Initialise a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{"sub": sub, "times": times}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all
	// of the 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through each of the pages.
	for _, page := range pages {
		// Extract the file name from the full file path and assign it to
		// the name variable
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create the empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		// ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		// if err != nil {
		// 	return nil, err
		// }

		// Do the same for partials.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
