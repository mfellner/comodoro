package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	"github.com/mfellner/comodoro/app"
	"github.com/mfellner/comodoro/db"
	"github.com/spf13/viper"
)

var dbFile *string

func init() {
	port := flag.Int("port", 8080, "Port to listen on")
	loglevel := flag.String("log", "info", "Loglevel (debug|info)")
	dbFile = flag.String("db", "/tmp/comodoro.db", "Path to the BoltDB file")
	flag.Parse()

	viper.SetEnvPrefix("app")

	viper.BindEnv("port")
	viper.BindEnv("loglevel")

	viper.SetDefault("loglevel", *loglevel)
	viper.Set("port", port)

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
	if err := db.Open(*dbFile, 0600); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := app.NewApp(&db)

	port := viper.GetInt("port")

	log.WithFields(log.Fields{
		"port": port,
	}).Info("comodoro started")

	log.Fatal(app.ListenAndServe(port))
}
