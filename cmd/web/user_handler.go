package main

import (
	"adilhaddad.net/agefice-docs/pkg/dao"
	"adilhaddad.net/agefice-docs/pkg/forms"
	model "adilhaddad.net/agefice-docs/pkg/models"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

func (app *application) SaveOrCheckLogin(w http.ResponseWriter, r *http.Request) {

	// parse r body as byte and then to player object

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var u model.Users
	{
	}
	type mess struct {
		Msg string
	}
	if err := json.Unmarshal(b, &u); err != nil {
		app.serverError(w, err)
		return
	}
	findUser, err := app.user.GetUserByEmail(r.Context(), u.Email)
	if err == nil && findUser.Email != "" {
		if app.comparePasswords(findUser.HashedPassword, []byte(u.HashedPassword)) {
			json.NewEncoder(w).Encode(&struct {
				Msg string `json:"msg"`
			}{Msg: "User identified"})
		} else {
			json.NewEncoder(w).Encode(&struct {
				Msg string `json:"msg"`
			}{Msg: "Mot de passe incorrect"})

		}
	} else {
		json.NewEncoder(w).Encode(&struct {
			Msg string `json:"msg"`
		}{Msg: err.Error()})
		hashedPass := app.hashAndSalt([]byte(u.HashedPassword))
		u.HashedPassword = []byte(hashedPass)

		result, _, err := app.user.AddUser(r.Context(), &u)
		if err != nil {
			json.NewEncoder(w).Encode("Erreur lors de la creation du user")
			returnError(w, r, err)
			return
		}
		writeJSON(r.Context(), w, result)

	}
}
func (app *application) signUpTemp(w http.ResponseWriter, r *http.Request) {
	app.renderPage(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{
		Form: forms.New(nil)}, r))
}
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {

	avatarUrl := app.manageAvatar(w, r)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.renderPage(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
	}

	var u model.Users
	u.NewUser(form.Get("name"), form.Get("email"), app.hashAndSalt([]byte(form.Get("password"))), avatarUrl)
	_, _, err = app.user.AddUser(r.Context(), &u)
	if err == dao.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.renderPage(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
	} else if err != nil {
		//returnError(w, r, err)
		form.Errors.Add("email", "Address is already in use")
		app.renderPage(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
		return
	}

	var f flash
	f.Code = "#36f439"
	f.Label = "INFO"
	f.Message = "Your signup was successful. Please log in."
	app.session.Put(r, "flash", f)

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
}
func (app *application) signInTemp(w http.ResponseWriter, r *http.Request) {
	app.renderPage(w, r, "signin.page.tmpl", app.addDefaultData(&templateData{
		Form: forms.New(nil)}, r))
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error // message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	result, err := app.user.GetUserByEmail(r.Context(), form.Get("email"))
	if err != nil {
		//returnError(w, r, err)
		form.Errors.Add("generic", "no matching record found")
		app.renderPage(w, r, "signin.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
	}
	// Check whether the hashed password and plain-text password provided match. // If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(result.HashedPassword, []byte(form.Get("password")))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			//returnError(w, r, dao.ErrInvalidCredentials)
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.renderPage(w, r, "signin.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
			return
		}
	}

	// Add the ID of the current user to the session, so that they are now 'logged // in'.
	app.session.Put(r, "authenticatedUserID", result.ID)
	var f flash
	f.Code = "#36f439"
	f.Label = "INFO"
	f.Message = "You've been login  successfully!"
	app.session.Put(r, "flash", f)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been // logged out.
	var f flash
	f.Code = "#36f439"
	f.Label = "INFO"
	f.Message = "You've been logged out successfully!"
	app.session.Put(r, "flash", f)
	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := int32(app.session.GetInt(r, "authenticatedUserID"))

	user, err := app.user.GetUser(r.Context(), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderPage(w, r, "profile.page.tmpl", &templateData{
		User: user,
	})
}

func (app *application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.renderPage(w, r, "password.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength("newPassword", 10)
	if form.Get("newPassword") != form.Get("newPasswordConfirmation") {
		form.Errors.Add("newPasswordConfirmation", "Passwords do not match")
	}

	if !form.Valid() {
		app.renderPage(w, r, "password.page.tmpl", &templateData{Form: form})
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")

	_, _, err = app.ChangePassword(r, userID, form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		//form.Errors.Add("currentPassword", "Current password is incorrect")
		//app.renderPage(w, r, "password.page.tmpl", &templateData{Form: form})
		returnError(w, r, err)
	}

	var f flash
	f.Code = "#36f439"
	f.Label = "INFO"
	f.Message = "Your password has been updated!"
	app.session.Put(r, "flash", f)

	http.Redirect(w, r, "/user/profile", 303)
}

func (app *application) GetProfile(w http.ResponseWriter, r *http.Request) {
	app.renderPage(w, r, "profile.page.tmpl", nil)
}
