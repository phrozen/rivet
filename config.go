package main

import (
	"flag"
    "os"
)

const (
    LIMIT = 1000
)

type Config struct {
    // Server
	Port          string
    // Database
    DatabasePath  string
}

func NewConfig() *Config {
    cfg := new(Config)
    flag.StringVar(&cfg.DatabasePath, "db", os.Getenv("RIVET_DATABASE_PATH"), "Path/folder to store/open the database files.")
	flag.StringVar(&cfg.Port, "port", os.Getenv("PORT"), "Port to listen.")
    flag.Parse()
    return cfg
}
    
