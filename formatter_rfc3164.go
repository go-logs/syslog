package syslog

import (
	"os"
	"fmt"
	"time"
)

const RFC3164_STRING_FORMAT_HEADER = "<%d>%s %s %s[%d]:"

// RFC3164Header 
type RFC3164Header struct {
	Hostname       string
	TimestampIsUTC bool
	Tag            string
}

// RFC3164 
type RFC3164 struct {
	Facility Priority
	Header   *RFC3164Header
}

// hostname 
func (h *RFC3164Header) hostname() {
	if h.Hostname == EMPTY_STRING {
		h.Hostname, _ = os.Hostname()
	}
	if len(h.Hostname) > HEADER_HOSTNAME_LENGTH {
		h.Hostname = truncateStartStr(h.Hostname, HEADER_HOSTNAME_LENGTH)
	}
}

// tag 
func (h *RFC3164Header) tag() {
	if h.Tag == EMPTY_STRING {
		h.Tag = os.Args[0]
	}
	if len(h.Tag) > HEADER_TAG_LENGTH {
		h.Tag = truncateStartStr(h.Tag, HEADER_TAG_LENGTH)
	}
}

// timestamp 
func (h *RFC3164Header) timestamp(tt time.Time) string {
	if h.TimestampIsUTC {
		tt = tt.UTC()
	}

	return tt.Format(time.Stamp)
}

// String 
func (h *RFC3164Header) String(priority Priority) string {
	return fmt.Sprintf(RFC3164_STRING_FORMAT_HEADER,
		priority, h.timestamp(time.Now()), h.Hostname, h.Tag, os.Getpid())
}

// Close 
func (h *RFC3164Header) Close() {
	if h != nil {
		h = nil
	}
}

// priority 
func (f *RFC3164) priority(severity Priority) Priority {
	return priority(BuildPriority(f.Facility, severity))
}

// string 
func (f *RFC3164) string(severity Priority) string {
	if f.Header == nil {
		f.Header = &RFC3164Header{}
	}

	f.Header.hostname()
	f.Header.tag()

	return f.Header.String(severity)
}

// String 
func (f *RFC3164) String(severity Priority, message string) string {
	if message == EMPTY_STRING {
		return f.string(severity)
	} else {
		return f.string(severity) + SPACE_STRING + message
	}
}

// Close 
func (f *RFC3164) Close() {
	if f != nil {
		if f.Header != nil {
			f.Header.Close()
			f.Header = nil
		}
		f = nil
	}
}

// RFC3164Formatter provides an RFC 3164 compliant message.
func RFC3164Formatter(p Priority, hostname, appName, tag, content string) string {
	r := &RFC3164{
		Facility: LOG_DAEMON,
		Header:   &RFC3164Header{
			Hostname:       hostname,
			Tag:            tag,
			TimestampIsUTC: true,
		},
	}
	return r.String(p, content)
}
