package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
	"testing/iotest"
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
	left := log.Colors.DarkGray + "["
	right := log.Colors.DarkGray + "]"
	for _, tc := range []struct {
		in   string
		want string
	}{
		{`{"level":"trace","msg":"foo"}`, left + log.Colors.Cyan + "VRB" + right + "> " +
			log.Colors.Cyan + "foo" + log.Colors.Reset + "\n"},
		{`{"level":"xyz","msg":"foo"}`, log.Colors.BrightRed + "?> " + log.Colors.Blue + "foo" + log.Colors.Reset + "\n"},
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

func TestScanLinesMax(t *testing.T) {
	input := "short\n" + strings.Repeat("x", 20) + "\nok\n" + strings.Repeat("y", 8) + "\nlast" + strings.Repeat("z", 30)
	expected := []string{"short", "xxxxxxxx", "ok", "yyyyyyyy", "lastzzzz"}
	readers := map[string]func(string) io.Reader{
		"whole":    func(s string) io.Reader { return strings.NewReader(s) },
		"one byte": func(s string) io.Reader { return iotest.OneByteReader(strings.NewReader(s)) },
	}
	for name, mk := range readers {
		scanner := bufio.NewScanner(mk(input))
		scanner.Buffer(make([]byte, 0, 4), 8)
		scanner.Split(scanLinesMax(8))
		var got []string
		for scanner.Scan() {
			got = append(got, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("%s: unexpected error %v", name, err)
		}
		if strings.Join(got, "|") != strings.Join(expected, "|") {
			t.Errorf("%s: got %q expected %q", name, got, expected)
		}
	}
}
