package main

import (
	"adilhaddad.net/agefice-docs/pkg/dao"
	m "adilhaddad.net/agefice-docs/pkg/models"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/gob"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/matryer/runner"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	session       *sessions.Session
	db            *gorm.DB
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	templateData  *templateData
	fi            *os.File
	fe            *os.File
	env           string
	recipe        interface {
		GetAllRecipes(ctx context.Context, page, pagesize int64, order string) (results []*m.Recipes, totalRows int64, err error)
		GetRecipe(ctx context.Context, argID int32) (record *m.Recipes, err error)
		AddRecipe(ctx context.Context, record *m.Recipes) (result *m.Recipes, RowsAffected int64, err error)
		UpdateRecipe(ctx context.Context, argID int32, updated *m.Recipes) (result *m.Recipes, RowsAffected int64, err error)
		DeleteRecipe(ctx context.Context, argID int32) (rowsAffected int64, err error)
	}
	user interface {
		GetAllUsers(ctx context.Context, page, pagesize int64, order string) (results []*m.Users, totalRows int64, err error)
		GetUser(ctx context.Context, argID int32) (record *m.Users, err error)
		AddUser(ctx context.Context, record *m.Users) (result *m.Users, RowsAffected int64, err error)
		UpdateUser(ctx context.Context, argID int32, updated *m.Users) (result *m.Users, RowsAffected int64, err error)
		DeleteUser(ctx context.Context, argID int32) (rowsAffected int64, err error)
		GetUserByEmail(ctx context.Context, argEmail string) (record *m.Users, err error)
	}
	serviceMail serviceMail
}

type flash struct {
	Code, Label, Message string
}

type config struct {
	addr      string
	staticDir string
}

type serviceMail struct {
	task       *runner.Task
	IsStarting bool
}

//type contextKey string

//var contextKeyUser = contextKey("isAuthenticated")

func main() {

	// get envirenement
	env := *flag.String("env", "DEV", "Envirenement")
	//env := *flag.String("env", "PROD", "Envirenement")

	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	gob.Register(flash{})

	session := sessions.New([]byte(*secret))
	session.Lifetime = 2 * time.Hour
	session.Secure = true

	//for windows OS
	fi, err := os.OpenFile("C:\\goLogs\\info.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	fe, err := os.OpenFile("C:\\goLogs\\error.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	//for lunix mac OS
	//fi, err := os.OpenFile("/tmp/info.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	//fe, err := os.OpenFile("/tmp/error.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
		fe.WriteString(err.Error())
	}
	defer fi.Close()
	defer fe.Close()
	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//infoLog = log.New(fi, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	//errorLog = log.New(fe, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	cfg := new(config)
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	//addr := flag.String("addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.addr, "addr", ":4001", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string.
	var dsn string
	if env == "DEV" {
		//mac
		//dsn = *flag.String("dsn", "root:sherine2011*@tcp(localhost:3306)/agefice_docs?parseTime=true", "MySQL data source name")
		//window
		dsn = *flag.String("dsn", "root:r00t@tcp(localhost:3306)/recipe?parseTime=true", "MySQL data source name")
	} else {
		dsn = *flag.String("dsn", "adil:toto@tcp(54.38.189.215:3306)/agefice_docs?parseTime=true", "MySQL data source name")
	}
	//dsn := flag.String("dsn", "root:r00t@tcp(localhost:3306)/agefice_docs?parseTime=true", "MySQL data source name")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openGormDB(dsn, fe)
	if err != nil {
		errorLog.Fatal(err)
		fe.WriteString(err.Error())
	}
	sqlDB, err := db.DB() // close connexion in gorm v2

	//We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	//defer db.Close()
	defer sqlDB.Close()

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		session:  session,
		db:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
		recipe:   &dao.RecipesDao{Db: db},
		user:     &dao.UsersDao{Db: db},

		fi:  fi,
		fe:  fe,
		env: env}

	// Initialize a new template cache...
	var templateCache map[string]*template.Template
	if env == "DEV" {
		templateCache, err = app.newTemplateCache("./ui/html/")
	} else {
		templateCache, err = app.newTemplateCache("/var/www/go/deploy/recipe/ui/html/")
	}
	// add templates cache to the application struct
	app.templateCache = templateCache

	//fmt.Printf("template cache => %s" ,templateCache)

	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want // the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256}}

	srv := &http.Server{
		Addr:     cfg.addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
		//TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server. IdleTimeout: time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %v ", cfg.addr)
	fi.WriteString("Starting server on : " + cfg.addr)
	if env == "DEV" {
		app.errorLog.Fatal(srv.ListenAndServe())
		//err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	} else {
		srv.TLSConfig = tlsConfig
		app.errorLog.Fatal(srv.ListenAndServeTLS("/var/www/go/deploy/agefice/tls/cert.pem", "/var/www/go/deploy/recipe/tls/key.pem"))
	}
	fe.WriteString(err.Error() + "\n")

}

// for a given DSN.
func openDB(dsn string, fe *os.File) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fe.WriteString(err.Error())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func openGormDB(dsn string, fe *os.File) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	//db, err := gorm.Open("mysql", dsn)
	//db.SingularTable(true) for singular name table
	if err != nil {
		fe.WriteString(err.Error())
		return nil, err
	}

	return db, nil

}
