package syslog

import (
	"fmt"
	"os"
	"time"
)

const (
	formatterDefaultString = "<%d> %s %s %s[%d]: %s"
	formatterUnixString    = "<%d>%s %s[%d]: %s"

	HEADER_PRIORITY_MIN = Priority(0)
	HEADER_PRIORITY_MAX = Priority((LOG_LOCAL7 & FacilityMask) | (LOG_DEBUG & SeverityMask))

	HEADER_HOSTNAME_LENGTH = 255

	HEADER_TAG_LENGTH = 32
)

// Formatter is a type of function that takes the consituent parts of a
// syslog message and returns a formatted string. A different Formatter is
// defined for each different syslog protocol we support.
type Formatter func(p Priority, hostname, appName, tag, content string) string
//type Formatter func(f interface{}) string

// DefaultFormatter is the original format supported by the Go syslog package,
// and is a non-compliant amalgamation of 3164 and 5424 that is intended to
// maximize compatibility.
func DefaultFormatter(p Priority, hostname, appName, tag, content string) string {
	timestamp := time.Now().Format(time.RFC3339)
	return fmt.Sprintf(formatterDefaultString,
		p, timestamp, hostname, tag, os.Getpid(), content)
}

// UnixFormatter omits the hostname, because it is only used locally.
func UnixFormatter(p Priority, hostname, appName, tag, content string) string {
	timestamp := time.Now().Format(time.Stamp)
	return fmt.Sprintf(formatterUnixString,
		p, timestamp, tag, os.Getpid(), content)
}
