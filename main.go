package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kaellybot/kaelly-twitter/application"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	initLog()
	initMetrics()
}

func initConfig() {
	viper.SetConfigFile(constants.ConfigFileName)

	for configName, defaultValue := range constants.GetDefaultConfigValues() {
		viper.SetDefault(configName, defaultValue)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Str(constants.LogFileName, constants.ConfigFileName).Msgf("Failed to read config file, continue...")
	}

	viper.AutomaticEnv()
}

func initLog() {
	gin.SetMode(gin.ReleaseMode)
	zerolog.SetGlobalLevel(constants.LogLevelFallback)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	log.Logger = log.With().Caller().Logger()

	logLevel, err := zerolog.ParseLevel(viper.GetString(constants.LogLevel))
	if err != nil {
		log.Warn().Err(err).Msgf("Log level not set, continue with %s...", constants.LogLevelFallback)
	} else {
		zerolog.SetGlobalLevel(logLevel)
		log.Debug().Msgf("Logger level set to '%s'", logLevel)
	}
}

func initMetrics() {
	go func() {
		log.Info().Msgf("Exposing Prometheus metrics...")
		http.Handle("/metrics", promhttp.Handler())

		server := &http.Server{
			Addr:              fmt.Sprintf(":%v", viper.GetInt(constants.MetricPort)),
			ReadHeaderTimeout: 0,
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Error().Err(err).Msgf("Cannot listen and serve Prometheus metrics")
		}
	}()
}

func main() {
	app := application.New()
	app.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Info().Msgf("%s v%s is now running. Press CTRL-C to exit.", constants.InternalName, constants.Version)
	<-sc

	log.Info().Msgf("Gracefully shutting down %s...", constants.InternalName)
	app.Shutdown()
}
