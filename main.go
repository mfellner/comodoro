package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mfellner/comodoro/app"
	"github.com/mfellner/comodoro/db"
	"github.com/spf13/viper"
)

var (
	dbFile        *string
	globalFlagset = flag.NewFlagSet("comodoro", flag.ExitOnError)

	globalFlags = struct {
		Port          int
		LogLevel      string
		DBFile        string
		FleetEndpoint string
	}{}
)

func init() {
	globalFlagset.IntVar(&globalFlags.Port,
		"port", 3030, "Port to listen on.")
	globalFlagset.StringVar(&globalFlags.LogLevel,
		"log", "info", "Log level (\"debug\", \"info\" or \"warn\").")
	globalFlagset.StringVar(&globalFlags.DBFile,
		"db", "/tmp/comodoro.db", "Path to the BoltDB file.")
	globalFlagset.StringVar(&globalFlags.FleetEndpoint,
		"fleet-endpoint", "unix:///var/run/fleet.sock", "Location of the fleet API.")
	globalFlagset.Parse(os.Args[1:])

	// Enviroment variables are uppercase and must start with APP_
	viper.SetEnvPrefix("app")
	viper.BindEnv("port")
	viper.BindEnv("loglevel")
	viper.BindEnv("fleetEndpoint")

	viper.SetDefault("port", globalFlags.Port)
	viper.SetDefault("loglevel", globalFlags.LogLevel)
	viper.SetDefault("fleetEndpoint", globalFlags.FleetEndpoint)

	switch viper.Get("loglevel") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	var db db.DB
	if err := db.Open(globalFlags.DBFile, 0600); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := app.NewApp(&db)

	port := viper.GetInt("port")

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Comodoro started")

	log.Fatal(app.ListenAndServe(port))
}
