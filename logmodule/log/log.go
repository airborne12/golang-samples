package log

import (
	"os"

	"github.com/op/go-logging"
)

var gLog = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.999} %{shortfile}:%{callpath} ▶ %{level} %{id}%{color:reset} %{message}`,
)

func init() {

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
}

func Get() *logging.Logger {
	return gLog
}
