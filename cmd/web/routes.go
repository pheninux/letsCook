package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.homeTemp))
	mux.Get("/recipe/create/form", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.recipeCreateTemp))
	mux.Post("/recipes/filter", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getRecipesByFilter))
	mux.Get("/get/recipe/:id", http.HandlerFunc(app.GetRecipe))
	mux.Get("/get/recipes/", http.HandlerFunc(app.GetAllRecipes))
	mux.Post("/create/recipe/", http.HandlerFunc(app.AddRecipe))
	mux.Post("/delete/recipe/", http.HandlerFunc(app.DeleteRecipe))
	mux.Post("/update/recipe/", http.HandlerFunc(app.UpdateRecipe))
	mux.Get("/", http.HandlerFunc(app.homeTemp))
	mux.Get("/user/signin", dynamicMiddleware.ThenFunc(app.signInTemp))
	mux.Post("/user/signin", dynamicMiddleware.ThenFunc(app.signIn))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signUpTemp))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signUp))
	mux.Get("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logout))
	mux.Get("/get/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.GetProfile))

	var fileServer http.Handler
	if app.env == "DEV" {
		fileServer = http.FileServer(http.Dir("./ui/static/"))
	} else {
		fileServer = http.FileServer(http.Dir("/var/www/go/deploy/recipe/ui/static/"))
	}

	//mux.Get("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
