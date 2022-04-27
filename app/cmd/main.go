package main

import (
	"fmt"
	"quote/api/app/api/health"
	"quote/api/app/api/quotes"
	"quote/api/app/config"
	"quote/api/app/internal/core/domain"
	"quote/api/app/internal/repositories"
	"quote/api/app/server"
	"quote/api/app/server/middleware"
	"quote/api/app/server/router"
	"quote/api/app/tools"
	"quote/api/app/tools/logger"
	"quote/api/app/utils/monitoring"

	"github.com/joho/godotenv"

	"strconv"
)

const (
	Env                 = "ENV"
	EnvLogLevel         = "LOG_LEVEL"
	EnvLogJsonOutput    = "LOG_JSON_OUTPUT"
	EnvPort             = "PORT"
	EnvDatabaseHost     = "DATABASE_HOST"
	EnvDatabase         = "DATABASE"
	EnvDatabaseUsername = "DATABASE_USERNAME"
	EnvDatabasePassword = "DATABASE_PASSWORD"
	EnvDatabasePort     = "DATABASE_PORT"
	EnvSentryDsn        = "SENTRY_DSN"
)

func main() {
	log := logger.NewLogger("quotes-api")

	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file. Using defaults")
	}

	environment := tools.EnvOr(Env, "development")
	logLevel := tools.EnvOr(EnvLogLevel, "debug")
	logJsonOutput := tools.EnvOr(EnvLogJsonOutput, "true")
	port := tools.EnvOr(EnvPort, "8080")
	host := tools.EnvOr(EnvDatabaseHost, "localhost")
	database := tools.EnvOr(EnvDatabase, "quotes-db")
	databaseUser := tools.EnvOr(EnvDatabaseUsername, "quotes-user")
	databasePass := tools.EnvOr(EnvDatabasePassword, "quotes-pass")
	databasePort := tools.EnvOr(EnvDatabasePort, "5432")
	sentryDsn := tools.EnvOr(EnvSentryDsn, "")

	enableJsonOutput, err := strconv.ParseBool(logJsonOutput)
	if err != nil {
		enableJsonOutput = true
	}

	configuration := config.Config{
		Env:  environment,
		Port: port,
		Logging: config.LoggingConfig{
			Level:            logLevel,
			EnableJSONOutput: enableJsonOutput,
		},
		Database: config.DatabaseConfig{
			Host:     host,
			Database: database,
			User:     databaseUser,
			Password: databasePass,
			Port:     databasePort,
		},
		Monitoring: config.Monitoring{
			Sentry: config.Sentry{
				Dsn: sentryDsn,
			},
		},
	}

	monitoring.InitializeMonitoring(configuration.Monitoring)

	srv := server.NewServer(&configuration)

	// middlewares for the server
	corsMiddleware := middleware.NewCORSMiddleware(configuration.CorsHeaders)
	loggingMiddleware := middleware.NewLoggingMiddleware(configuration.Logging)
	recoveryMiddleware := middleware.NewRecoveryMiddleware()
	monitoringMiddleware := middleware.NewMonitoringMiddleware()

	// use middlewares
	srv.UseMiddleware(recoveryMiddleware)
	srv.UseMiddleware(monitoringMiddleware)
	srv.UseMiddleware(loggingMiddleware)
	srv.UseMiddleware(corsMiddleware)

	repository := repositories.NewRepository(configuration.Database)
	quotesService := domain.NewQuotesUseCase(repository.GetQuotesRepo())

	// setup routers
	routers := []router.Router{
		quotes.NewQuotesRouter(quotesService),
		health.NewHealthRouter(),
	}

	// initialize routers
	srv.InitRouter(routers...)

	appServer := srv.CreateServer()

	// start & run the server
	err = appServer.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		_, msg := fmt.Printf("Failed to start Server %s", err)
		log.Error(msg)
		panic(msg)
	}
}
