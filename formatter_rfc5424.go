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

	RFC5424_DATA_ID_NAME_TIME_QUALITY = "timeQuality"
	RFC5424_DATA_ID_NAME_TIME_QUALITY_TZ_KNOWN      = "tzKnown"
	RFC5424_DATA_ID_NAME_TIME_QUALITY_IS_SYNCED     = "isSynced"
	RFC5424_DATA_ID_NAME_TIME_QUALITY_SYNC_ACCURACY = "syncAccuracy"

	RFC5424_DATA_ID_NAME_ORIGIN       = "origin"
	RFC5424_DATA_ID_NAME_ORIGIN_IP            = "ip"
	RFC5424_DATA_ID_NAME_ORIGIN_ENTERPRISE_ID = "enterpriseId"
	RFC5424_DATA_ID_NAME_ORIGIN_SOFTWARE      = "software"
	RFC5424_DATA_ID_NAME_ORIGIN_SW_VERSION    = "swVersion"

	RFC5424_DATA_ID_NAME_META         = "meta"
	RFC5424_DATA_ID_NAME_META_SEQUENCE_ID = "sequenceId"
	RFC5424_DATA_ID_NAME_META_SYS_UP_TIME = "sysUpTime"
	RFC5424_DATA_ID_NAME_META_LANGUAGE    = "language"

	RFC5424_DATA_ID_META_SEQUENCE_ID_MIN = 0
	RFC5424_DATA_ID_META_SEQUENCE_ID_MAX = int32(2147483647)

	RFC5424_DATA_ID_META_SYS_UP_TIME_MIN = 0
	RFC5424_DATA_ID_META_SYS_UP_TIME_MAX = 9
)

const (
	RFC5424_TIMESTAMP_LEVEL_NONE  = "none"
	RFC5424_TIMESTAMP_LEVEL_MILLI = "milli"
	RFC5424_TIMESTAMP_LEVEL_MICRO = "micro"

	RFC5424_TIMESTAMP_FORMAT_MILLI = "2006-01-02T15:04:05.999Z07:00"
	RFC5424_TIMESTAMP_FORMAT_MICRO = "2006-01-02T15:04:05.999999Z07:00"
)

// RFC5424DataIDTimeQuality 
type RFC5424DataIDTimeQuality struct {
	TzKnown      bool
	IsSynced     bool
	SyncAccuracy int32
}

// RFC5424DataIDOrigin 
type RFC5424DataIDOrigin struct {
	IP           []string
	EnterpriseId string
	Software     string
	SwVersion    string
}

// RFC5424DataIDMeta 
type RFC5424DataIDMeta struct {
	SequenceId int32	// 2147483647
	SysUpTime  int		// 0-9
	Language   string
}

// RFC5424DataIDs 
type RFC5424DataIDs struct {
	TimeQuality *RFC5424DataIDTimeQuality
	Origin      *RFC5424DataIDOrigin
	Meta        *RFC5424DataIDMeta
}

// RFC5424DataParams 
type RFC5424DataParams map[string]interface{}

// RFC5424Data 
type RFC5424Data struct {
	ID     string
	Params RFC5424DataParams
}

// RFC5424StructuredData 
type RFC5424StructuredData struct {
	Elements []*RFC5424Data
}

// RFC5424Header 
type RFC5424Header struct {
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

// timeQualityString 
func (s *RFC5424DataIDs) timeQualityString() string {
	var str string

	if s != nil && s.TimeQuality != nil {
		str = RFC5424_DATA_BEGIN_LINE + RFC5424_DATA_ID_NAME_TIME_QUALITY +
			fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_TIME_QUALITY_TZ_KNOWN, boolInt(s.TimeQuality.TzKnown)) +
			fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_TIME_QUALITY_IS_SYNCED, boolInt(s.TimeQuality.IsSynced))
		if s.TimeQuality.IsSynced && s.TimeQuality.SyncAccuracy > 0 {
			str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_TIME_QUALITY_SYNC_ACCURACY, s.TimeQuality.SyncAccuracy)
		}
		str = str + RFC5424_DATA_END_LINE
	}

	return str
}

// originString 
func (s *RFC5424DataIDs) originString() string {
	var str string

	if s != nil && s.Origin != nil {
		str = RFC5424_DATA_BEGIN_LINE + RFC5424_DATA_ID_NAME_ORIGIN
		for i := 0; i < len(s.Origin.IP); i++ {
			str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_ORIGIN_IP, s.Origin.IP[i])
		}
		str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_ORIGIN_ENTERPRISE_ID, s.Origin.EnterpriseId) +
			fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_ORIGIN_SOFTWARE, s.Origin.Software) +
			fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_ORIGIN_SW_VERSION, s.Origin.SwVersion) +
			RFC5424_DATA_END_LINE
	}

	return str
}

// metaString 
func (s *RFC5424DataIDs) metaString() string {
	var str string

	if s != nil && s.Meta != nil {
		str = RFC5424_DATA_BEGIN_LINE + RFC5424_DATA_ID_NAME_META

		if s.Meta.SequenceId > RFC5424_DATA_ID_META_SEQUENCE_ID_MAX {
			s.Meta.SequenceId = RFC5424_DATA_ID_META_SEQUENCE_ID_MAX
		} else if s.Meta.SequenceId < RFC5424_DATA_ID_META_SEQUENCE_ID_MIN {
			s.Meta.SequenceId = RFC5424_DATA_ID_META_SEQUENCE_ID_MIN
		}
		str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
			RFC5424_DATA_ID_NAME_META_SEQUENCE_ID, s.Meta.SequenceId)

		if s.Meta.SysUpTime > RFC5424_DATA_ID_META_SYS_UP_TIME_MAX {
			s.Meta.SysUpTime = RFC5424_DATA_ID_META_SYS_UP_TIME_MAX
		} else if s.Meta.SysUpTime < RFC5424_DATA_ID_META_SYS_UP_TIME_MIN {
			s.Meta.SysUpTime = RFC5424_DATA_ID_META_SYS_UP_TIME_MIN
		}
		str = str + fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
			RFC5424_DATA_ID_NAME_META_SYS_UP_TIME, s.Meta.SysUpTime) +
			fmt.Sprintf(RFC5424_DATA_PARAM_FORMAT_STRING,
				RFC5424_DATA_ID_NAME_META_LANGUAGE, s.Meta.Language)
	}

	return str
}

// String 
func (s *RFC5424DataIDs) String() string {
	var str string

	if s != nil {
		str = str + s.timeQualityString() + s.originString() + s.metaString()
		if str != EMPTY_STRING {
			str = SPACE_STRING + str
		}
	}

	return str
}

// Close 
func (s *RFC5424DataIDs) Close() {
	if s != nil {
		s.TimeQuality = nil
		s.Origin = nil
		s.Meta = nil
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

// hostname 
func (h *RFC5424Header) hostname() {
	var err error

	if h.Hostname == EMPTY_STRING {
		h.Hostname, err = os.Hostname()
		if err != nil {
			h.Hostname = RFC5424_EMPTY_VALUE
			err = nil
		}
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
func (h *RFC5424Header) String(priority Priority) string {
	return fmt.Sprintf(RFC5424_STRING_FORMAT_HEADER,
		priority, RFC5424_VERSION, h.timestamp(time.Now()), h.Hostname, h.AppName, os.Getpid(), h.MessageID)
}

// Close 
func (h *RFC5424Header) Close() {
	if h != nil {
		h = nil
	}
}

// priority 
func (f *RFC5424) priority(severity Priority) Priority {
	return priority(BuildPriority(f.Facility, severity))
}

// structuredDataIDs 
func (f *RFC5424) structuredDataIDs() {
	if f.StructuredDataIDs == nil {
		f.StructuredDataIDs = &RFC5424DataIDs{}
	}
}

// structuredData 
func (f *RFC5424) structuredData() {
	if f.StructuredData == nil {
		f.StructuredData = &RFC5424StructuredData{}
	}
}

// header 
func (f *RFC5424) header() {
	if f.Header == nil {
		f.Header = &RFC5424Header{}
	}
}

// string 
func (f *RFC5424) string(severity Priority) string {
	f.header()
	f.Header.hostname()
	f.Header.appName()
	f.Header.messageID()

	return f.Header.String(f.priority(severity)) + SPACE_STRING + f.StructuredData.String() + f.StructuredDataIDs.String()
}

// String
func (f *RFC5424) String(severity Priority, message string) string {
	if message == EMPTY_STRING {
		return f.string(severity)
	} else {
		return f.string(severity) + SPACE_STRING + message
	}
}

// SetHostname 
func (f *RFC5424) SetHostname(hn string) {
	f.header()
	f.Header.Hostname = hn
}

// SetAppName 
func (f *RFC5424) SetAppName(an string) {
	f.header()
	f.Header.AppName = an
}

// SetFacility 
func (f *RFC5424) SetTimestampIsUTC(t bool) {
	f.header()
	f.Header.TimestampIsUTC = t
}

// SetTimestampLevel 
func (f *RFC5424) SetTimestampLevel(l string) {
	f.header()
	f.Header.TimestampLevel = l
}

// SetTag 
func (f *RFC5424) SetTag(t string) {
	f.header()
	f.Header.MessageID = t
}

// SetFacility 
func (f *RFC5424) SetFacility(facility Priority) {
	f.Facility = facility
}

// AddStructuredData 
func (f *RFC5424) AddStructuredData(data *RFC5424Data) {
	f.structuredData()
	f.StructuredData.Elements = append(f.StructuredData.Elements, data)
}

// SetStructuredData 
func (f *RFC5424) SetStructuredData(data *RFC5424Data) {
	if f.StructuredData != nil {
		f.StructuredData.Close()
	}
	f.AddStructuredData(data)
}

// SetStructuredDataIDs 
func (f *RFC5424) SetStructuredDataIDs(data *RFC5424DataIDs) {
	if f.StructuredDataIDs != nil {
		f.StructuredDataIDs.Close()
	}
	f.StructuredDataIDs = data
}

// SetStructuredDataIDTimeQuality 
func (f *RFC5424) SetStructuredDataIDTimeQuality(data *RFC5424DataIDTimeQuality) {
	f.structuredDataIDs()
	f.StructuredDataIDs.TimeQuality = data
}

// SetStructuredDataIDOrigin 
func (f *RFC5424) SetStructuredDataIDOrigin(data *RFC5424DataIDOrigin) {
	f.structuredDataIDs()
	f.StructuredDataIDs.Origin = data
}

// SetStructuredDataIDMeta 
func (f *RFC5424) SetStructuredDataIDMeta(data *RFC5424DataIDMeta) {
	f.structuredDataIDs()
	f.StructuredDataIDs.Meta = data
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
