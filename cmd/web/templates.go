package main

import (
	"adilhaddad.net/agefice-docs/pkg/forms"
	model "adilhaddad.net/agefice-docs/pkg/models"
	"html/template"
	"net/url"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear       int
	FormatDate        func(t time.Time) string
	Flash             flash
	AuthenticatedUser *model.Users
	IsAuthenticated   bool
	CSRFToken         string
	Form              *forms.Form
	FormData          url.Values
	FormErrors        map[string]string
	User              *model.Users
	PagedResults      *PagedResults
	TempTitle         string
	Recipes           []*model.Recipes
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humainDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("01/02/2006 15:04")
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func (app *application) getUser(id int) model.Users {
	u, _ := app.user.GetUser(nil, int32(id))
	return *u
}

// function qui renvoie un struct avec la moyenne des notes , nombre des avis
func rateAverageCalc(avis []model.Avis) interface{} {

	var star = 5
	var sumRate int    // moyenne des avis
	var s []int        //array data for start
	var f bool = false // boolean si en a fini de parcourir la moyenne de des notes

	for _, v := range avis {
		sumRate += v.Rate
	}

	for i := 0; i < star; i++ {
		if !f && sumRate > 0 {
			for j := 0; j < sumRate/len(avis); j++ {
				s = append(s, 1)
			}
			f = true
			i = sumRate / len(avis)
			if i == 5 {
				break
			}
		}
		s = append(s, 0)
	}

	return &struct {
		Avg       int
		AvisCount int
		Rates     []int
	}{
		Avg:       sumRate / len(avis),
		AvisCount: len(avis),
		Rates:     s,
	}
}

// function qui renvoie un struct avec la moyenne des notes , nombre des avis
func rateCalc(avis model.Avis) interface{} {

	var star = 5
	var s []int        //array data for start
	var f bool = false // boolean si en a fini de parcourir la moyenne de des notes

	for i := 0; i < star; i++ {
		if !f && avis.Rate > 0 {
			for j := 0; j < avis.Rate; j++ {
				s = append(s, 1)
			}
			f = true
			i = avis.Rate
			if i == 5 {
				break
			}
		}
		s = append(s, 0)
	}

	return &struct {
		Rates []int
	}{
		Rates: s,
	}
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functionsOld = template.FuncMap{
	"humainDate": humainDate,
	"formatDate": formatDate,
	"rateArg":    rateAverageCalc,
}

func (app *application) functions() template.FuncMap {
	return template.FuncMap{
		"humainDate":  humainDate,
		"formatDate":  formatDate,
		"rateAvgCalc": rateAverageCalc,
		"rateCalc":    rateCalc,
		"getUser":     app.getUser,
	}
}

func (app *application) newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {

		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {

		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(app.functions()).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts

	}
	return cache, nil
}
