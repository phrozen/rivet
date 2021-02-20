package main

//"fmt"

// App ...
type App struct {
	System *BoltDB
	Store  map[string]*BoltDB
	Config *Config
}

// NewApp ...
func NewApp(cfg *Config) *App {
	app := new(App)

	if cfg == nil {
		app.Config = NewConfig()
	} else {
		app.Config = cfg
	}

	var err error
	app.System, err = NewBoltDB(app.Config.DatabasePath + "system.db")
	if err != nil {
		panic(err)
	}

	systemBuckets := []string{"admin", "user", "session", "stat", "log"}
	for _, bucket := range systemBuckets {
		err := app.System.CreateBucketIfNotExist(bucket)
		if err != nil {
			panic(err)
		}
	}

	app.Store = make(map[string]*BoltDB, 0)
	users := app.System.All("user", "", 0)

	for _, u := range users {
		app.Store[u], err = NewBoltDB(app.Config.DatabasePath + "/store/" + u + ".db")
		if err != nil {
			panic(err)
		}
		err := app.Store[u].CreateBucketIfNotExist("store")
		if err != nil {
			panic(err)
		}
	}
	return app
}
