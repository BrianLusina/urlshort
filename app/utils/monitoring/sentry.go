package monitoring

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type Sentry struct {
	Dsn string
}

func NewSentry(dsn string) *Sentry {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	return &Sentry{Dsn: dsn}
}

func (s *Sentry) CaptureException(err error) {
	sentry.CaptureException(err)
}

func (s *Sentry) CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}
