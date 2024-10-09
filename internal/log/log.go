package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}
