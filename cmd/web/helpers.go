package main

import (
	"adilhaddad.net/agefice-docs/pkg/dao"
	model "adilhaddad.net/agefice-docs/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// PagedResults results for pages GetAll results.
type PagedResults struct {
	Page         int64       `json:"page"`
	PageSize     int64       `json:"page_size"`
	Data         interface{} `json:"data"`
	TotalRecords int64       `json:"total_records"`
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	trace := fmt.Sprintf("%s\n%s\n", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	app.fe.WriteString(fmt.Sprintf("%s\n%s\n", err.Error(), debug.Stack()))
	fmt.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) BodyParser(r *http.Request) []byte {
	body, _ := ioutil.ReadAll(r.Body)
	return body
}

func (app *application) ToJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=UTF8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		app.serverError(w, err)
		return
	}

}
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()

	if app.session.Exists(r, "flash") {
		td.Flash = app.session.Pop(r, "flash").(flash)
	}
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) renderPage(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
	}
	// Initialize a new buffer.
	buff := new(bytes.Buffer)
	err := ts.Execute(buff, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buff.WriteTo(w)

}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, file, tmpl string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.

	ts, ok := app.templateCache[file]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", file))
	}
	// Initialize a new buffer.
	buff := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buff, tmpl, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buff.WriteTo(w)

}

func (app *application) addCurrentYear(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) validateForm(r *http.Request) map[string]string {
	errors := make(map[string]string)

	name := r.PostForm.Get("nom")
	prenom := r.PostForm.Get("prenom")

	if strings.TrimSpace(name) == "" {
		errors["name"] = "le nom ne dois pas être vide"
	}
	if strings.TrimSpace(prenom) == "" {
		errors["prenom"] = "le prenom ne dois pas être vide"
	} else if utf8.RuneCountInString(prenom) > 25 {
		errors["prenom"] = "le prenom depasse le nombre de caractaires authorisé (max 10)"
	}

	return errors
}

func (app *application) createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func (app *application) FormatDate(date time.Time) string {
	const LAYOUT = "2006-01-02"
	t, _ := time.Parse(LAYOUT, date.String())
	return t.String()
}

/*** hash and salt password **/
func (app *application) hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

/*** compare password ***/
func (app *application) comparePasswords(hashedPwd []byte, plainPwd []byte) bool {
	var b bool
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		b = false
	} else {
		b = true
	}
	return b
}

/*func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}*/

// Return true if the current request is from authenticated user, otherwise return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}

// HTTPNocacheContent will set the headers for content type along with no caching.
func HTTPNocacheContent(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", content)
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// HTTPNocacheJSON will set the headers on an http Response for a text/json content type along with no cache.
func HTTPNocacheJSON(w http.ResponseWriter) {
	HTTPNocacheContent(w, "text/json")
}

// SendJSON will return take a value and serialize it to json and return the http response.
func SendJSON(w http.ResponseWriter, r *http.Request, code int, val interface{}) {
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	bytes, err := json.MarshalIndent(val, "", "")
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	_, _ = w.Write(bytes)
}

// InternalServerError will return an error to the client, sending 500 error code to the client with generic string
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Error(w, "Internal server error", 500)
}

// AddHeadersHandler will take a map of string/string and use it to set the key and value as the header name and value respectively.
func AddHeadersHandler(addHeaders map[string]string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for key, value := range addHeaders {
			w.Header().Set(key, value)
		}

		h.ServeHTTP(w, r)
	})
}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// IsPrivateSubnet - check to see if this ip is in a private subnet
func IsPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

// GetIPAddress will take a http request and check headers if it has been proxied to extract what the server believes to be the client ip address.
func GetIPAddress(r *http.Request) string {
	ip := ""
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip = strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}

	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

// FormatRequest generates ascii representation of a request
func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func (app *application) ChangePassword(r *http.Request, id int, currentPassword, newPassword string) (result *model.Users, RowsAffected int64, err error) {
	u, err := app.user.GetUser(r.Context(), int32(id))
	if err != nil {
		return nil, 0, err
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, 0, dao.ErrInvalidCredentials
		} else {
			return nil, 0, err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return nil, 0, err
	}
	u.HashedPassword = newHashedPassword
	result, row, err := app.user.UpdateUser(r.Context(), int32(id), u)
	return result, row, err

}

func initializeContext(r *http.Request) (ctx context.Context) {
	if ContextInitializer != nil {
		ctx = ContextInitializer(r)
	} else {
		ctx = r.Context()
	}
	return ctx
}

func ValidateRequest(ctx context.Context, r *http.Request, table string, action model.Action) error {
	if RequestValidator != nil {
		return RequestValidator(ctx, r, table, action)
	}

	return nil
}

type RequestValidatorFunc func(ctx context.Context, r *http.Request, table string, action model.Action) error

var RequestValidator RequestValidatorFunc

type ContextInitializerFunc func(r *http.Request) (ctx context.Context)

var ContextInitializer ContextInitializerFunc

func readInt(r *http.Request, param string, v int64) (int64, error) {
	p := r.FormValue(param)
	if p == "" {
		return v, nil
	}

	return strconv.ParseInt(p, 10, 64)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	data, _ := json.MarshalIndent(v, "", " ")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func writeRowsAffected(w http.ResponseWriter, rowsAffected int64) {
	data, _ := json.Marshal(rowsAffected)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func readJSON(r *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, v)
}

func returnError(w http.ResponseWriter, r *http.Request, err error) {
	status := 0
	switch err {
	case dao.ErrMismatchedHashAndPassword:
		status = http.StatusInternalServerError
	case dao.ErrInvalidCredentials:
		status = http.StatusInternalServerError
	case dao.ErrDuplicateEmail:
		status = http.StatusInternalServerError
	case dao.ErrNoRecordFound:
		status = http.StatusInternalServerError
	case dao.ErrUnableToMarshalJSON:
		status = http.StatusInternalServerError
	case dao.ErrUpdateFailed:
		status = http.StatusInternalServerError
	case dao.ErrInsertFailed:
		status = http.StatusInternalServerError
	case dao.ErrDeleteFailed:
		status = http.StatusInternalServerError
	case dao.ErrBadParams:
		status = http.StatusBadRequest
	case dao.ErrScanColumns:
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	SendJSON(w, r, er.Code, er)
}

func parseUint8(r *http.Request, key string) (uint8, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 8)
	if err != nil {
		return uint8(id), err
	}
	return uint8(id), err
}
func parseUint16(r *http.Request, key string) (uint16, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		return uint16(id), err
	}
	return uint16(id), err
}
func parseUint32(r *http.Request, key string) (uint32, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return uint32(id), err
	}
	return uint32(id), err
}
func parseUint64(r *http.Request, key string) (uint64, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return uint64(id), err
	}
	return uint64(id), err
}
func parseInt(r *http.Request, key string) (int, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(id), err
}
func parseInt8(r *http.Request, key string) (int8, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 8)
	if err != nil {
		return -1, err
	}
	return int8(id), err
}
func parseInt16(r *http.Request, key string) (int16, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		return -1, err
	}
	return int16(id), err
}
func parseInt32(r *http.Request, key string) (int32, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return -1, err
	}
	return int32(id), err
}
func parseInt64(r *http.Request, key string) (int64, error) {
	idStr := r.URL.Query().Get(key)
	id, err := strconv.ParseInt(idStr, 10, 54)
	if err != nil {
		return -1, err
	}
	return id, err
}
func parseString(r *http.Request, key string) (string, error) {
	idStr := r.URL.Query().Get(key)
	return idStr, nil
}
func parseUUID(r *http.Request, key string) (string, error) {
	idStr := r.URL.Query().Get(key)
	return idStr, nil
}

func (app *application) manageAvatar(w http.ResponseWriter, r *http.Request) string {

	var avatarUrl string
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(1 << 2)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		//fmt.Println("Error Retrieving the File")
		//fmt.Println(err)
		//return
		avatarUrl = `C:\\Utilisateurs\\a706836\\go\\src\\recepe-manager\\avatar-temp\\defaultAvatar.png`
	} else {
		avatarUrl = `C:\\Utilisateurs\\a706836\\go\\src\\recepe-manager\\avatar-temp\\` + filepath.Base(handler.Filename)
		defer file.Close()

		//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		//fmt.Printf("File Size: %+v\n", handler.Size)
		//fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create file
		//dst, err := os.Create(handler.Filename)
		////////////////////////////////////////////////////
		dst, err := os.Create(filepath.Join(avatarUrl))
		if err != nil {
			fmt.Println(err)
			app.serverError(w, err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			fmt.Println(err)
			app.serverError(w, err)
		}
	}

	return avatarUrl
}
