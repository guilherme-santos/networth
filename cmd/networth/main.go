package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/foodora/go-ranger/fdhttp"
	"github.com/guilherme-santos/networth/http"
	"github.com/guilherme-santos/networth/mongodb"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var config struct {
	Port       string `default:"80"`
	MongoDBURL string `required:"true"`
	LogFormat  string `envconfig:"NETWORTH_LOG_FORMAT" default:"json"`
	LogLevel   string `envconfig:"NETWORTH_LOG_LEVEL" default:"info"`
}

func main() {
	var log = logrus.New()

	err := envconfig.Process("networth", &config)
	if err != nil {
		log.WithError(err).Fatal("Unable to load environment variables")
	}

	switch config.LogFormat {
	case "text":
		log.Formatter = &logrus.TextFormatter{}
	default:
		log.Formatter = &logrus.JSONFormatter{}
	}

	var lvl logrus.Level
	if lvl, err = logrus.ParseLevel(config.LogLevel); err != nil {
		log.WithError(err).
			WithField("log_level", config.LogLevel).
			Warn("Unable to identify log level")
		lvl = logrus.InfoLevel
	}
	log.Level = lvl

	fdhttp.SetLogger(log)

	lm := fdhttp.NewLogMiddleware()
	lm.SetLogger(log)

	storage := mongodb.NewStorage(config.MongoDBURL)

	router := fdhttp.NewRouter()
	router.Use(lm.Middleware())
	router.Register(http.NewHandler(storage))

	srv := fdhttp.NewServer(config.Port)

	var errChan chan error
	go func() {
		errChan <- srv.Start(router)
	}()

	stopSignal := make(chan os.Signal, 2)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	// block until receive a SIGTERM or server.Start return
	select {
	case <-stopSignal:
		err := srv.Stop()
		if err != nil {
			log.WithError(err).Fatal("Unable to stop gracefully")
		}
	case err := <-errChan:
		log.WithError(err).Fatalf("Unable to run http server")
	}

	log.Print("Application stopped succesfuly!")
}
