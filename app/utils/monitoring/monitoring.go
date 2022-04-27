package monitoring

import "quote/api/app/config"

type Monitoring interface {
	CaptureException(err error)
	CaptureMessage(message string)
}

func InitializeMonitoring(config config.Monitoring) {
	NewSentry(config.Sentry.Dsn)
}
