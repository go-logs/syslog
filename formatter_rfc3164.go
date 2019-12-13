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
	h.Hostname = BuildHostname(h.Hostname)
}

// tag 
func (h *RFC3164Header) tag() {
	h.Tag = BuildTag(h.Tag)
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
	return BuildPriority(f.Facility, severity)
}

// string 
func (f *RFC3164) headerString(severity Priority) string {
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
		return f.headerString(severity)
	} else {
		return f.headerString(severity) + SPACE_STRING + message
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
			TimestampIsUTC: false,
		},
	}
	return r.String(p, content)
}
