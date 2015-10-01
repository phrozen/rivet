package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"os"
)

type App struct {
	dbs *bolt.DB
	dba *bolt.DB

	dbsPath string
	dbaPath string
	port    string
	cli     bool
}

func NewApp() *App {
	app := new(App)

	flag.StringVar(&app.dbsPath, "dbs", os.Getenv("RIVET_DBS"), "Path to the store database.")
	flag.StringVar(&app.dbaPath, "dba", os.Getenv("RIVET_DBA"), "Path to the system (admin) database.")
	flag.StringVar(&app.port, "port", os.Getenv("PORT"), "Port to listen.")
	flag.BoolVar(&app.cli, "cli", false, "Command line interface.")

	return app
}

func (app *App) Initialize() {
	flag.Parse()

	dbs, err := bolt.Open(app.dbsPath, 0600, nil)
	if err != nil {
		panic("Store database path not found! (Defaults to environment variable RIVET_DBS or set with flag -dbs)")
	}
	app.dbs = dbs

	dba, err := bolt.Open(app.dbaPath, 0600, nil)
	if err != nil {
		panic("System (admin) database path not found! (Defaults to environment variable RIVET_DBA or set with flag -dba)")
	}
	app.dba = dba

	systemBuckets := []string{"admin", "user", "session", "stat", "log"}
	for _, bucket := range systemBuckets {
		app.dba.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				panic(err)
			}
			return nil
		})
	}

	// Create mockup user
	err = app.SetUserPassword("test", "test")
	if err != nil {
		panic(err)
	}
	app.dbs.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			panic(err)
		}
		return nil
	})
}
