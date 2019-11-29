package syslog

import (
	"fmt"
	"os"
	"time"
)

const appNameMaxLength = 48  // limit to 48 chars as per RFC5424

const (
	formatterDefaultString     = "<%d> %s %s %s[%d]: %s"
	formatterUnixString        = "<%d>%s %s[%d]: %s"
	formatterUnixRFC3164String = "<%d>%s %s %s[%d]: %s"
	formatterUnixRFC5424String = "<%d>%d %s %s %s %d %s - %s"
)

// Formatter is a type of function that takes the consituent parts of a
// syslog message and returns a formatted string. A different Formatter is
// defined for each different syslog protocol we support.
type Formatter func(p Priority, hostname, appName, tag, content string) string

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

// RFC3164Formatter provides an RFC 3164 compliant message.
func RFC3164Formatter(p Priority, hostname, appName, tag, content string) string {
	timestamp := time.Now().Format(time.Stamp)
	return fmt.Sprintf(formatterUnixRFC3164String,
		p, timestamp, hostname, tag, os.Getpid(), content)
}

// if string's length is greater than max, then use the last part
func truncateStartStr(s string, max int) string {
	if (len(s) > max) {
		return s[len(s) - max:]
	}
	return s
}

// RFC5424Formatter provides an RFC 5424 compliant message.
func RFC5424Formatter(p Priority, hostname, appName, tag, content string) string {
	timestamp := time.Now().Format(time.RFC3339)
	if len(appName) == 0 {
		appName = truncateStartStr(os.Args[0], appNameMaxLength)
	}
	return fmt.Sprintf(formatterUnixRFC5424String,
		p, 1, timestamp, hostname, appName, os.Getpid(), tag, content)
}
