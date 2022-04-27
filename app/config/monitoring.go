package config

type Sentry struct {
	Dsn string
}

type Monitoring struct {
	Sentry
}
