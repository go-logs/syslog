package syslog

import (
	"fmt"
	"os"
	"time"
)

const (
	RFC5424_VERSION = 1

	RFC5424_EMPTY_VALUE = "-"

	RFC5424_STRING_FORMAT_HEADER = "<%d>%d %s %s %s %d %s"

	RFC5424_HEADER_VERSION_MIN = 0
	RFC5424_HEADER_VERSION_MAX = 2

	RFC5424_HEADER_APP_NAME_LENGTH = 48

	RFC5424_DATA_PARAM_FORMAT_STRING = " %s=\"%v\""
	RFC5424_DATA_BEGIN_LINE          = "["
	RFC5424_DATA_END_LINE            = "]"
)

const (
	RFC5424_TIMESTAMP_LEVEL_NONE  = "none"
	RFC5424_TIMESTAMP_LEVEL_MILLI = "milli"
	RFC5424_TIMESTAMP_LEVEL_MICRO = "micro"

	RFC5424_TIMESTAMP_FORMAT_MILLI = "2006-01-02T15:04:05.999Z07:00"
	RFC5424_TIMESTAMP_FORMAT_MICRO = "2006-01-02T15:04:05.999999Z07:00"
)

// RFC5424DataIDsTimeQuality 
type RFC5424DataIDsTimeQuality struct {
	tzKnown      bool
	isSynced     bool
	syncAccuracy int32
}

// RFC5424DataIDsOrigin 
type RFC5424DataIDsOrigin struct {
	ip           string
	enterpriseId string
	software     string
	swVersion    string
}

// RFC5424DataIDsMeta 
type RFC5424DataIDsMeta struct {
	sequenceId string
	sysUpTime  string
	language   string
}

// RFC5424DataIDs 
type RFC5424DataIDs struct {
	timeQuality *RFC5424DataIDsTimeQuality
	origin      *RFC5424DataIDsOrigin
	meta        *RFC5424DataIDsMeta
}

// RFC5424DataParams 
type RFC5424DataParams map[string]interface{}

// RFC5424Data 
type RFC5424Data struct {
	ID       string
	Params RFC5424DataParams
}

// RFC5424StructuredData 
type RFC5424StructuredData struct {
	Elements []*RFC5424Data
}

// RFC5424Header 
type RFC5424Header struct {
	Priority       Priority
	Version        int
	Hostname       string
	AppName        string
	TimestampIsUTC bool
	TimestampLevel string
	MessageID      string
}

// RFC5424 
type RFC5424 struct {
	Facility          Priority
	Header            *RFC5424Header
	StructuredData    *RFC5424StructuredData
	StructuredDataIDs *RFC5424DataIDs
}

// if string's length is greater than max, then use the last part
func truncateStartStr(s string, max int) string {
	if (len(s) > max) {
		return s[len(s) - max:]
	}
	return s
}

// boolInt 
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// String.
func (s *RFC5424DataIDs) String() string {
	var str string

	if s != nil {
		str = SPACE_STRING + str
	}

	return str
}

// Close.
func (s *RFC5424DataIDs) Close() {
	if s != nil {
		s.timeQuality = nil
		s.origin = nil
		s.meta = nil
		s = nil
	}
}

// String 
func (d *RFC5424Data) String() string {
	str := RFC5424_DATA_BEGIN_LINE + d.ID

	for key, value := range d.Params {
		str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING, key, value)
	}

	return str + RFC5424_DATA_END_LINE
}

// Close 
func (d *RFC5424Data) Close() {
	if d != nil {
		d.Params = nil
		d = nil
	}
}

// String 
func (s *RFC5424StructuredData) String() string {
	var str string

	if s != nil {
		for i := 0; i < len(s.Elements); i++ {
			if s.Elements[i].ID != EMPTY_STRING && len(s.Elements[i].Params) > 0 {
				str = str + s.Elements[i].String()
			}
		}
	}

	if str == EMPTY_STRING {
		str = RFC5424_EMPTY_VALUE
	}

	return str
}

// Close 
func (s *RFC5424StructuredData) Close() {
	if s != nil {
		for i := 0; i < len(s.Elements); i++ {
			s.Elements[i].Close()
		}
		s = nil
	}
}

// priority 
func (h *RFC5424Header) priority() {
	if h.Priority < HEADER_PRIORITY_MIN {
		h.Priority = HEADER_PRIORITY_MIN
	} else if h.Priority > HEADER_PRIORITY_MAX {
		h.Priority = HEADER_PRIORITY_MAX
	}
}

// setPriority 
func (h *RFC5424Header) setPriority(facility Priority, severity Priority) {
	h.Priority = (facility & FacilityMask) | (severity & SeverityMask)
}

// hostname 
func (h *RFC5424Header) hostname() {
	if h.Hostname == EMPTY_STRING {
		h.Hostname, _ = os.Hostname()
	}
	if len(h.Hostname) > HEADER_HOSTNAME_LENGTH {
		h.Hostname = truncateStartStr(h.Hostname, HEADER_HOSTNAME_LENGTH)
	}
}

// appName 
func (h *RFC5424Header) appName() {
	if h.AppName == EMPTY_STRING {
		h.AppName = os.Args[0]
	}
	if len(h.AppName) > RFC5424_HEADER_APP_NAME_LENGTH {
		h.AppName = truncateStartStr(h.AppName, RFC5424_HEADER_APP_NAME_LENGTH)
	}
}

// messageID 
func (h *RFC5424Header) messageID() {
	if h.MessageID == EMPTY_STRING {
		h.MessageID = RFC5424_EMPTY_VALUE
	}
	if h.MessageID != RFC5424_EMPTY_VALUE && len(h.MessageID) > HEADER_TAG_LENGTH {
		h.MessageID = truncateStartStr(h.MessageID, HEADER_TAG_LENGTH)
	}
}

// timestamp 
func (h *RFC5424Header) timestamp(tt time.Time) string {
	if h.TimestampIsUTC {
		tt = tt.UTC()
	}

	switch h.TimestampLevel {
	case RFC5424_TIMESTAMP_LEVEL_MILLI:
		return tt.Format(RFC5424_TIMESTAMP_FORMAT_MILLI)
	case RFC5424_TIMESTAMP_LEVEL_MICRO:
		return tt.Format(RFC5424_TIMESTAMP_FORMAT_MICRO)
	default:
		return tt.Format(time.RFC3339)
	}
}

// String 
func (h *RFC5424Header) String() string {
	return fmt.Sprintf(RFC5424_STRING_FORMAT_HEADER,
		h.Priority, RFC5424_VERSION, h.timestamp(time.Now()), h.Hostname, h.AppName, os.Getpid(), h.MessageID)
}

// Close 
func (h *RFC5424Header) Close() {
	if h != nil {
		h = nil
	}
}

// string 
func (f *RFC5424) string() string {
	if f.Header == nil {
		f.Header = &RFC5424Header{}
	}

	f.Header.priority()
	f.Header.hostname()
	f.Header.appName()
	f.Header.messageID()

	return f.Header.String() + SPACE_STRING + f.StructuredData.String() + f.StructuredDataIDs.String()
}

// String
func (f *RFC5424) String(severity Priority, message string) string {
	f.Header.setPriority(f.Facility, severity)

	if message == EMPTY_STRING {
		return f.string()
	} else {
		return f.string() + SPACE_STRING + message
	}
}

// Close 
func (f *RFC5424) Close() {
	if f != nil {
		if f.Header != nil {
			f.Header.Close()
			f.Header = nil
		}
		if f.StructuredData != nil {
			f.StructuredData.Close()
			f.StructuredData = nil
		}
		if f.StructuredDataIDs != nil {
			f.StructuredDataIDs.Close()
			f.StructuredDataIDs = nil
		}
		f = nil
	}
}

// RFC5424Formatter provides an RFC 5424 compliant message
func RFC5424Formatter(p Priority, hostname, appName, tag, content string) string {
	r := &RFC5424{
		Facility: LOG_DAEMON,
		Header:   &RFC5424Header{
			Hostname:       hostname,
			AppName:        appName,
			MessageID:      tag,
			TimestampIsUTC: true,
			TimestampLevel: RFC5424_TIMESTAMP_LEVEL_MILLI,
		},
//		StructuredData:    &RFC5424StructuredData{[]*RFC5424Data{&RFC5424Data{ID: "fields", Params: RFC5424DataParams{"arg": 123, "arg2": "param"}}}},
//		StructuredDataIDs: &RFC5424DataIDs{},
	}
	return r.String(p, content)
}
