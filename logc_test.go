package main

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"fortio.org/log"
)

func TestGetAttributes(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{`{"msg":"foo"}`, ""},
		{`{"msg":"foo","attr1":"val\"with quote"}`, `: "attr1":"val\"with quote"`},
		{`{"msg":"foo\\"}`, ``},
		{`{"msg":"foo\\","a1":"v1"}`, `: "a1":"v1"`},
		{`{"msg":"hello \"world\"","a1":"v\n1,"a2":"v2"}`, `: "a1":"v\n1,"a2":"v2"`},
	} {
		got := GetAttributes(tc.in)
		if got != tc.want {
			t.Errorf("GetAttributes(%s) = '%s', want '%s'", tc.in, got, tc.want)
		}
	}
}

func TestLevels(t *testing.T) {
	var zeroTime time.Time
	log.Config.ForceColor = true
	log.SetColorMode()
	for _, tc := range []struct {
		in   string
		want string
	}{
		{`{"level":"trace","msg":"foo"}`, log.Colors.Cyan + "Verb" + log.Colors.DarkGray + ">" +
			log.Colors.Cyan + "foo" + log.Colors.Reset + "\n"},
		{`{"level":"xyz","msg":"foo"}`, log.Colors.Blue + "? foo" + log.Colors.Reset + "\n"},
	} {
		buf := &strings.Builder{}
		w := bufio.NewWriter(buf)
		ProcessLogLine(w, &zeroTime, []byte(tc.in))
		w.Flush()
		got := buf.String()
		if got != tc.want {
			t.Errorf("LevelToColor(%s) = '%s', want '%s'", tc.in, got, tc.want)
		}
	}
}
