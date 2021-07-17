package main

import (
	"adilhaddad.net/agefice-docs/pkg/dao"
	model "adilhaddad.net/agefice-docs/pkg/models"
	"fmt"
	"net/http"

	"github.com/guregu/null"
)

var (
	_ = null.Bool{}
)

func (app *application) homeTemp(w http.ResponseWriter, r *http.Request) {

	ctx := initializeContext(r)
	page, err := readInt(r, "page", 0)
	if err != nil || page < 0 {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	pagesize, err := readInt(r, "pagesize", 20)
	if err != nil || pagesize <= 0 {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	order := r.FormValue("order")

	if err := ValidateRequest(ctx, r, "recipes", model.RetrieveMany); err != nil {
		returnError(w, r, err)
		return
	}

	records, totalRows, err := app.recipe.GetAllRecipes(ctx, page, pagesize, order)
	if err != nil {
		returnError(w, r, err)
		return
	}

	result := &PagedResults{Page: page, PageSize: pagesize, Data: records, TotalRecords: totalRows}
	app.renderPage(w, r, "home.page.tmpl", &templateData{PagedResults: result, TempTitle: "List of recipes"})
}

func (app *application) recipeCreateTemp(w http.ResponseWriter, r *http.Request) {

}
func (app *application) recipeShowTemp(w http.ResponseWriter, r *http.Request) {

}

// GetAllRecipes is a function to get a slice of record(s) from recipes table in the recipe database
// @Summary Get list of Recipes
// @Tags Recipes
// @Description GetAllRecipes is a handler to get a slice of record(s) from recipes table in the recipe database
// @Accept  json
// @Produce  json
// @Param   page     query    int     false        "page requested (defaults to 0)"
// @Param   pagesize query    int     false        "number of records in a page  (defaults to 20)"
// @Param   order    query    string  false        "db sort order column"
// @Success 200 {object} api.PagedResults{data=[]model.Recipes}
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /recipes [get]
// http "http://localhost:8080/recipes?page=0&pagesize=20" X-Api-User:user123
func (app *application) GetAllRecipes(w http.ResponseWriter, r *http.Request) {

	ctx := initializeContext(r)
	page, err := readInt(r, "page", 0)
	if err != nil || page < 0 {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	pagesize, err := readInt(r, "pagesize", 20)
	if err != nil || pagesize <= 0 {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	order := r.FormValue("order")

	if err := ValidateRequest(ctx, r, "recipes", model.RetrieveMany); err != nil {
		returnError(w, r, err)
		return
	}

	records, totalRows, err := app.recipe.GetAllRecipes(ctx, page, pagesize, order)
	if err != nil {
		returnError(w, r, err)
		return
	}

	result := &PagedResults{Page: page, PageSize: pagesize, Data: records, TotalRecords: totalRows}
	writeJSON(ctx, w, result)
}

// GetRecipes is a function to get a single record from the recipes table in the recipe database
// @Summary Get record from table Recipes by  argID
// @Tags Recipes
// @ID argID
// @Description GetRecipes is a function to get a single record from the recipes table in the recipe database
// @Accept  json
// @Produce  json
// @Param  argID path int true "id"
// @Success 200 {object} model.Recipes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError "ErrNotFound, db record for id not found - returns NotFound HTTP 404 not found error"
// @Router /recipes/{argID} [get]
// http "http://localhost:8080/recipes/1" X-Api-User:user123
func (app *application) GetRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := initializeContext(r)

	argID, err := parseInt32(r, "argID")
	if err != nil {
		returnError(w, r, err)
		return
	}

	if err := ValidateRequest(ctx, r, "recipes", model.RetrieveOne); err != nil {
		returnError(w, r, err)
		return
	}

	record, err := app.recipe.GetRecipe(ctx, argID)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(ctx, w, record)
}

// AddRecipes add to add a single record to recipes table in the recipe database
// @Summary Add an record to recipes table
// @Description add to add a single record to recipes table in the recipe database
// @Tags Recipes
// @Accept  json
// @Produce  json
// @Param Recipes body model.Recipes true "Add Recipes"
// @Success 200 {object} model.Recipes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /recipes [post]
// echo '{"id": 72,"title": "ZhiGIPpVVaiOijXLPJLZjTpDV","descri": "XhuZluqbxeFsjYenPUcwBxfbD","obs": "ABeyALQcYrvMchYGtSVcUCjxM","categorie": "afReyIgsvMWfFZCgbteGtVNKu","preparation": 37,"typ": "wFhkuIMqESHdJolCbaAXGNLon","cuisson": 15,"repos": 11,"lvl": "jcHbbDHXQoJDCHKsiBkCjtMqh","nbr_pers": 85,"cout": 0.64038724,"share": "aklXTnWNyAPFeLJirAhEXWeUd","valide": 22,"id_users": 97}' | http POST "http://localhost:8080/recipes" X-Api-User:user123
func (app *application) AddRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := initializeContext(r)
	recipes := &model.Recipes{}

	if err := readJSON(r, recipes); err != nil {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	if err := recipes.BeforeSave(app.db); err != nil {
		returnError(w, r, dao.ErrBadParams)
	}

	recipes.Prepare()

	if err := recipes.Validate(model.Create); err != nil {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	if err := ValidateRequest(ctx, r, "recipes", model.Create); err != nil {
		returnError(w, r, err)
		return
	}

	var err error
	recipes, _, err = app.recipe.AddRecipe(ctx, recipes)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(ctx, w, recipes)
}

// UpdateRecipes Update a single record from recipes table in the recipe database
// @Summary Update an record in table recipes
// @Description Update a single record from recipes table in the recipe database
// @Tags Recipes
// @Accept  json
// @Produce  json
// @Param  argID path int true "id"
// @Param  Recipes body model.Recipes true "Update Recipes record"
// @Success 200 {object} model.Recipes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /recipes/{argID} [put]
// echo '{"id": 72,"title": "ZhiGIPpVVaiOijXLPJLZjTpDV","descri": "XhuZluqbxeFsjYenPUcwBxfbD","obs": "ABeyALQcYrvMchYGtSVcUCjxM","categorie": "afReyIgsvMWfFZCgbteGtVNKu","preparation": 37,"typ": "wFhkuIMqESHdJolCbaAXGNLon","cuisson": 15,"repos": 11,"lvl": "jcHbbDHXQoJDCHKsiBkCjtMqh","nbr_pers": 85,"cout": 0.64038724,"share": "aklXTnWNyAPFeLJirAhEXWeUd","valide": 22,"id_users": 97}' | http PUT "http://localhost:8080/recipes/1"  X-Api-User:user123
func (app *application) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := initializeContext(r)

	argID, err := parseInt32(r, "argID")
	if err != nil {
		returnError(w, r, err)
		return
	}

	recipes := &model.Recipes{}
	if err := readJSON(r, recipes); err != nil {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	if err := recipes.BeforeSave(app.db); err != nil {
		returnError(w, r, dao.ErrBadParams)
	}

	recipes.Prepare()

	if err := recipes.Validate(model.Update); err != nil {
		returnError(w, r, dao.ErrBadParams)
		return
	}

	if err := ValidateRequest(ctx, r, "recipes", model.Update); err != nil {
		returnError(w, r, err)
		return
	}

	recipes, _, err = app.recipe.UpdateRecipe(ctx, argID, recipes)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(ctx, w, recipes)
}

// DeleteRecipes Delete a single record from recipes table in the recipe database
// @Summary Delete a record from recipes
// @Description Delete a single record from recipes table in the recipe database
// @Tags Recipes
// @Accept  json
// @Produce  json
// @Param  argID path int true "id"
// @Success 204 {object} model.Recipes
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /recipes/{argID} [delete]
// http DELETE "http://localhost:8080/recipes/1" X-Api-User:user123
func (app *application) DeleteRecipe(w http.ResponseWriter, r *http.Request) {

	ctx := initializeContext(r)

	argID, err := parseInt32(r, "argID")
	if err != nil {
		returnError(w, r, err)
		return
	}

	if err := ValidateRequest(ctx, r, "recipes", model.Delete); err != nil {
		returnError(w, r, err)
		return
	}

	rowsAffected, err := app.recipe.DeleteRecipe(ctx, argID)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeRowsAffected(w, rowsAffected)
}

func (app *application) getRecipesByFilter(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
	}

	fmt.Println(r.PostForm.Get("recipe-title"))
	fmt.Println(r.PostForm.Get("preparation-time"))
	fmt.Println(r.PostForm.Get("numbers-covers"))
	fmt.Println(r.PostForm.Get("level"))
	fmt.Println(r.PostForm.Get("cost"))

	records, totalRows, err := app.recipe.GetAllRecipes(r.Context(), 0, 20, "")
	if err != nil {
		returnError(w, r, err)
		return
	}

	result := &PagedResults{Page: 0, PageSize: 20, Data: records, TotalRecords: totalRows}
	app.renderTemplate(w, r, "recipes.layout.tmpl", "recipes", &templateData{PagedResults: result, TempTitle: "List of recipes"})
}
