package main

import (
	ormsql "github.com/Pioneersltd/DevDailyDigest/v1/pkg/models/mysql"
	"flag"
	"log"
	"net/http"
	"os"
)

type App struct {
	infoLog *log.Logger
	warnLog *log.Logger
	errLog  *log.Logger
	model   *ormsql.UserModel
	bearer  *string
	origin  *string
}

func main() {

	addr := flag.String("addr", ":8081", "Address that you want the golang api to run on")
	sqlUser := flag.String("dbuser", "", "Mysql user")
	pass := flag.String("pass", "", "Mysql pass")
	dbSchema := flag.String("dbtable", "", "Database table")
	bearerToken := flag.String("token", "", "Bearer token")
	origin := flag.String("origin", "", "Host header")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	warnLogger := log.New(os.Stdout, "WARN:\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := OpenDB(*dbSchema, *sqlUser, *pass, "parseTime=true")

	if err != nil {
		errLogger.Fatal(err)
	}
	defer db.Close()

	app := App{
		infoLog: infoLogger,
		warnLog: warnLogger,
		errLog:  errLogger,
		model:   &ormsql.UserModel{DB: db},
		bearer:  bearerToken,
		origin:  origin,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLogger,
		Handler:  app.routes(),
	}

	infoLogger.Printf("Starting up HTTP server on %v", *addr)
	if err := srv.ListenAndServe(); err != nil {
		errLogger.Fatal(err)
	}

}
