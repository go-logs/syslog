package syslog

import (
	"os"
	"fmt"
	"time"
	"strings"
	"testing"
)

func TestDefaultFormatter(t *testing.T) {
	out := DefaultFormatter(LOG_ERR, "hostname", "appName", "tag", "content")
	expected := fmt.Sprintf("<%d> %s %s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.RFC3339), "hostname", "tag", os.Getpid(), "content")
	if out != expected {
		t.Errorf("expected %v got %v", expected, out)
	}
}

func TestUnixFormatter(t *testing.T) {
	out := UnixFormatter(LOG_ERR, "hostname", "appName", "tag", "content")
	expected := fmt.Sprintf("<%d>%s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.Stamp), "tag", os.Getpid(), "content")
	if out != expected {
		t.Errorf("expected %v got %v", expected, out)
	}
}

func TestRFC5424Formatter(t *testing.T) {
//	out := RFC5424Formatter(LOG_ERR, "hostname", truncateStartStr(os.Args[0], RFC5424_HEADER_APP_NAME_LENGTH), "tag", "content")
//	expected := fmt.Sprintf("<%d>%d %s %s %s %d %s - - %s",
//		LOG_ERR, 1, time.Now().Format(time.RFC3339), "hostname", truncateStartStr(os.Args[0], RFC5424_HEADER_APP_NAME_LENGTH),
//			os.Getpid(), "tag", "content")
//	if out != expected {
//		t.Errorf("expected %v got %v", expected, out)
//	}

//	out = RFC5424Formatter(LOG_ERR, "hostname", "", "tag", "content")
//	expected = fmt.Sprintf("<%d>%d %s %s %s %d %s - - %s",
//		LOG_ERR, 1, time.Now().Format(time.RFC3339), "hostname", truncateStartStr(os.Args[0], RFC5424_HEADER_APP_NAME_LENGTH),
//			os.Getpid(), "tag", "content")
//	if out != expected {
//		t.Errorf("Should be same and expected %v got %v", expected, out)
//	}

	out := RFC5424Formatter(LOG_ERR, "hostname", "appName", "tag", "content")
	expected := fmt.Sprintf("<%d>%d %s %s %s %d %s - - %s",
		LOG_ERR, 1, time.Now().Format(time.RFC3339), "hostname", truncateStartStr(os.Args[0], RFC5424_HEADER_APP_NAME_LENGTH),
			os.Getpid(), "tag", "content")
	if out == expected {
		t.Errorf("Should be differ and expected %v got %v", expected, out)
	}
}

func TestTruncateStartStr(t *testing.T) {
	out := truncateStartStr("abcde", 3)
	if strings.Compare(out, "cde" ) != 0 {
		t.Errorf("expected \"cde\" got %v", out)
	}
	out = truncateStartStr("abcde", 5)
	if strings.Compare(out, "abcde" ) != 0 {
		t.Errorf("expected \"abcde\" got %v", out)
	}
}

func TestPriorityLimits(t *testing.T) {
	priority := priorityLimits(HEADER_PRIORITY_MIN - 1)
	if priority != HEADER_PRIORITY_MIN {
		t.Errorf("Min priority should be equal HEADER_PRIORITY_MIN(%d) but got %d", priority, HEADER_PRIORITY_MIN)
	}
	priority = priorityLimits(HEADER_PRIORITY_MAX + 1)
	if priority != HEADER_PRIORITY_MAX {
		t.Errorf("Max priority should be equal HEADER_PRIORITY_MAX(%d) but got %d", priority, HEADER_PRIORITY_MAX)
	}
}

func TestBuildPriority(t *testing.T) {
	if BuildPriority(LOG_KERN, LOG_EMERG) != HEADER_PRIORITY_MIN {
		t.Errorf("Priority should be equal HEADER_PRIORITY_MIN(%d)", HEADER_PRIORITY_MIN)
	}

	if BuildPriority(LOG_LOCAL7, LOG_DEBUG) != HEADER_PRIORITY_MAX {
		t.Errorf("Priority should be equal HEADER_PRIORITY_MAX(%d)", HEADER_PRIORITY_MAX)
	}
}

func TestBuildHostname(t *testing.T) {
	hn, _ := os.Hostname()
	currentHn := BuildHostname(EMPTY_STRING)

	if hn != currentHn {
		t.Errorf("Hostname should be same and expected %s\ngot %s", hn, currentHn)
	}

	hn = strings.Repeat("hostname_test", 20)
	currentHn = BuildHostname(hn)
	if currentHn == hn || len(currentHn) == len(hn) {
		t.Errorf("Hostname length should be not more HEADER_HOSTNAME_LENGTH(%d)", HEADER_HOSTNAME_LENGTH)
	}
}

func TestBuildTag(t *testing.T) {
	tag := os.Args[0]
	currentTag := BuildTag(EMPTY_STRING)

	if !(currentTag == tag || strings.Contains(tag, currentTag)) {
		t.Errorf("Tag should be equal %s for %s", tag, currentTag)
	}

	tag = strings.Repeat("tag_test", 10)
	if currentTag == tag || len(currentTag) == len(tag) {
		t.Errorf("Tag length should be not more HEADER_TAG_LENGTH(%d)", HEADER_TAG_LENGTH)
	}
}
