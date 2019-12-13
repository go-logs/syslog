package syslog

import (
	"os"
	"fmt"
	"time"
	"strings"
	"testing"
)

func TestRFC3164HeaderHostname(t *testing.T) {
	h := &RFC3164Header{}
	h.hostname()
	hn, _ := os.Hostname()

	if h.Hostname != hn {
		t.Errorf("RFC3164 hostname should be equal %s for RFC3164Header %s", hn, h.Hostname)
	}

	hn = strings.Repeat("hostname_test", 20)
	h.Hostname = hn
	h.hostname()

	if h.Hostname == hn || len(h.Hostname) == len(hn) {
		t.Errorf("RFC3164 hostname length should be not more HEADER_HOSTNAME_LENGTH(%d)", HEADER_HOSTNAME_LENGTH)
	}
	h.Close()
}

func TestRFC3164HeaderTag(t *testing.T) {
	h := &RFC3164Header{}
	h.tag()
	tag := os.Args[0]

	if !(h.Tag == tag || strings.Contains(tag, h.Tag)) {
		t.Errorf("RFC3164 tag should be equal %s for RFC3164Header %s", tag, h.Tag)
	}

	tag = strings.Repeat("tag_test", 10)
	h.Tag = tag
	h.tag()

	if h.Tag == tag || len(h.Tag) == len(tag) {
		t.Errorf("RFC3164 tag length should be not more HEADER_TAG_LENGTH(%d)", HEADER_TAG_LENGTH)
	}
	h.Close()
}

func TestRFC3164HeaderTimestamp(t *testing.T) {
	h := &RFC3164Header{}
	tt := time.Now()
	ts := tt.Format(time.Stamp)

	if h.timestamp(tt) != ts {
		t.Errorf("RFC3164 timestamp should be equal %s", ts)
	}

	h.TimestampIsUTC = true
	ts = tt.UTC().Format(time.Stamp)
	if h.timestamp(tt) != ts {
		t.Errorf("RFC3164 timestamp should be equal %s as UTC", ts)
	}
	h.Close()
}

func TestRFC3164HeaderString(t *testing.T) {
	h := &RFC3164Header{}
	h.Hostname = "hostname"
	h.Tag = "tag"

	str := fmt.Sprintf(RFC3164_STRING_FORMAT_HEADER,
		HEADER_PRIORITY_MIN, h.timestamp(time.Now()), h.Hostname, h.Tag, os.Getpid())
	if h.String(HEADER_PRIORITY_MIN) != str {
		t.Errorf("RFC3164 string should be equal %s", str)
	}
	h.Close()
}

func TestRFC3164HeaderClose(t *testing.T) {
	h := &RFC3164Header{
		Hostname:       "test_hostname",
		TimestampIsUTC: true,
		Tag:            "test_tag",
	}
	h.Close()

	if h == nil {
		t.Error("Should be RFC3164Header destroyed")
	}
}

func TestRFC3164Priority(t *testing.T) {
	f := &RFC3164{}

	f.Facility = LOG_KERN
	if f.priority(LOG_EMERG) != HEADER_PRIORITY_MIN {
		t.Errorf("RFC3164 priority should be equal HEADER_PRIORITY_MIN(%d)", HEADER_PRIORITY_MIN)
	}

	f.Facility = LOG_LOCAL7
	if f.priority(LOG_DEBUG) != HEADER_PRIORITY_MAX {
		t.Errorf("RFC3164 priority should be equal HEADER_PRIORITY_MAX(%d)", HEADER_PRIORITY_MAX)
	}
	f.Close()
}

func TestRFC3164MainHeaderString(t *testing.T) {
	f := &RFC3164{}
	strDefault := f.headerString(LOG_KERN)
	strCurrent := fmt.Sprintf(RFC3164_STRING_FORMAT_HEADER,
		f.priority(LOG_KERN), f.Header.timestamp(time.Now()), f.Header.Hostname, f.Header.Tag, os.Getpid())

	if strCurrent != strDefault {
		t.Errorf("RFC3164 Header strings and should be equls:\n%s\nand\n%s\n", strCurrent, strDefault)
	}
	f.Close()
}

func TestRFC3164String(t *testing.T) {
	f := &RFC3164{}
	str := f.String(LOG_KERN, EMPTY_STRING)
	strHeader := f.headerString(LOG_KERN)
	if str != strHeader {
		t.Errorf("RFC3164 string without message and should be equls:\n%s\nand\n%s\n", str, strHeader)
	}

	msg := "MESSAGE"
	str = f.String(LOG_KERN, msg)
	strHeader = f.headerString(LOG_KERN)
	strCurrent := strHeader + SPACE_STRING + msg
	if strCurrent != str {
		t.Errorf("RFC3164 string with message should be %s\n", str)
	}
	f.Close()
}

func TestRFC3164Close(t *testing.T) {
	f := &RFC3164{
		HEADER_PRIORITY_MAX,
		&RFC3164Header{
			Hostname:       "test_hostname",
			TimestampIsUTC: true,
			Tag:            "test_tag",
		},
	}
	f.Close()

	if f == nil {
		t.Error("Should be RFC3164 destroyed")
	}
}

func TestRFC3164SetHostname(t *testing.T) {
	f := &RFC3164{}
	hn := "test_hostname"
	f.SetHostname(hn)

	if f.Header.Hostname != hn {
		t.Errorf("Set RFC3164 hostname should be equal %s for RFC3164Header %s", hn, f.Header.Hostname)
	}
	f.Close()
}

func TestRFC3164SetTimestampIsUTC(t *testing.T) {
	f := &RFC3164{}
	timestampIsUTC := true
	f.SetTimestampIsUTC(timestampIsUTC)

	if f.Header.TimestampIsUTC != timestampIsUTC {
		t.Errorf("Set RFC3164 UTC timestamp should be equal %t for RFC3164Header %t", timestampIsUTC, f.Header.TimestampIsUTC)
	}
	f.Close()
}

func TestRFC3164SetTag(t *testing.T) {
	f := &RFC3164{}
	tag := "test_tag"
	f.SetTag(tag)

	if f.Header.Tag != tag {
		t.Errorf("Set RFC3164 tag should be equal %s for RFC3164Header %s", tag, f.Header.Tag)
	}
	f.Close()
}

func TestRFC3164SetFacility(t *testing.T) {
	f := &RFC3164{}
	f.SetFacility(LOG_SYSLOG)

	if f.Facility != LOG_SYSLOG {
		t.Errorf("Set RFC3164 facility should be equal %d for RFC3164Header %d", LOG_SYSLOG, f.Facility)
	}
	f.Close()
}

func TestRFC3164Formatter(t *testing.T) {
	out := RFC3164Formatter(LOG_ERR, "hostname", "appName", "tag", "content")
	expected := fmt.Sprintf("<%d>%s %s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.Stamp), "hostname", "tag", os.Getpid(), "content")

	if out != expected {
		t.Errorf("RFC3164 Formatter message  should be equls but expected\n%v\ngot\n%v", expected, out)
	}
}
