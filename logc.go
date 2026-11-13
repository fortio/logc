package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"fortio.org/cli"
	"fortio.org/log"
)

// LevelToColor takes a json string level and returns the 1 letter short and color to use for it.
func LevelToColor(levelStr string) (string, string) {
	level, found := log.JSONStringLevelToLevel[levelStr]
	if !found {
		log.Critf("Bug/Unknown level %q", levelStr)
		return log.Colors.BrightRed + "?", log.Colors.Blue
	}
	log.Debugf("level %q -> %d", levelStr, level)
	return log.ColorLevelToStr(level), log.LevelToColor[level]
}

// GetAttributes returns the remaining/additional attributes after the `msg`, if any, for text output.
// faster than reparsing as a map.
func GetAttributes(line string) string {
	idx1 := strings.Index(line, `"msg":"`)
	if idx1 < 0 {
		log.Critf("Bug line without msg tag %q", line)
		return ""
	}
	i := idx1 + 7
	for {
		idx2 := strings.Index(line[i:], `"`)
		if idx2 < 0 {
			log.Critf("Bug line without close quote %q", line)
			return ""
		}
		i += idx2 + 1
		// not an escaped quote inside msg (\") but \\" is ok (see test)
		if line[i-2] != '\\' || line[i-3] == '\\' {
			break
		}
	}
	end := len(line) - 1
	if i == end {
		log.Debugf("no attributes for %q", line)
		return ""
	}
	log.Debugf("found attributes at %d/%d for %q", i, end, line)
	return ": " + line[i+1:end] // better more efficient way?
}

func main() {
	noColorFlag := flag.Bool("no-color", false, "Do not colorize output (same as -logger-no-color)")
	cli.ArgsHelp = " < log.json\nor for instance\n\tfortio server 2>&1 | logc\n" +
		"to convert JSON fortio logger lines from stdin to (ansi) colorized text"
	cli.Main()
	// read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 4096), maxLineSize)
	scanner.Split(scanLinesMax(maxLineSize))
	prevDate := time.UnixMilli(0)
	if *noColorFlag {
		log.Debugf("Disabling color - was %v", log.Color)
		log.Config.ConsoleColor = false
		log.SetColorMode()
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		ProcessLogLine(os.Stdout, &prevDate, line)
	}
	if err := scanner.Err(); err != nil {
		log.Errf("Error reading stdin: %v", err)
	}
}

// maxLineSize is the longest line we process, longer ones are truncated.
const maxLineSize = 64 * 1024

// scanLinesMax is like bufio.ScanLines but truncates lines longer than max (which must be
// at least the scanner's max token size) instead of failing with bufio.ErrTooLong,
// which would permanently stop the scanner (eg NULs from a truncated log under tail -f).
func scanLinesMax(maxLen int) bufio.SplitFunc {
	skipping := false // discarding the remainder of a truncated line
	return func(data []byte, atEOF bool) (int, []byte, error) {
		advance, token, err := bufio.ScanLines(data, atEOF)
		if skipping {
			if token == nil && err == nil { // still no end of line, drop what we have
				return len(data), nil, nil
			}
			skipping = false // found the end of the line, resume with the next one
			return advance, nil, err
		}
		if token == nil && err == nil && len(data) >= maxLen {
			log.Warnf("Line longer than %d bytes, truncating", maxLen)
			skipping = true
			return len(data), data[:maxLen], nil
		}
		return advance, token, err
	}
}

func ProcessLogLine(w io.Writer, prevDate *time.Time, line []byte) {
	// json deserialize
	e := log.JSONEntry{}
	err := json.Unmarshal(line, &e)
	if err != nil {
		log.LogVf("Error unmarshalling %q: %v", string(line), err)
		fmt.Printf("! %s\n", string(line))
		return
	}
	tsStr := ""
	// uppercase single letter level + color extraction
	lvl, color := LevelToColor(e.Level)
	if e.TS != 0 {
		ts := e.Time()
		// Each time the day changes we print a header
		if ts.YearDay() != prevDate.YearDay() {
			fmt.Printf("#### %s ####\n", ts.Format("2006-01-02"))
			*prevDate = ts
		}
		// Use full microseconds resolution unlike the log built in color version which stops at millis.
		tsStr = ts.Format(log.Colors.DarkGray + "15:04:05.000000 ")
	}
	if e.R > 0 {
		tsStr += fmt.Sprintf(log.Colors.Gray+"[%d] ", e.R)
	}
	fileLine := ""
	if e.Line != 0 {
		fileLine = fmt.Sprintf("%s:%d> ", e.File, e.Line)
	} else {
		lvl += ">"
	}
	// Msg can be multi line.
	fmt.Fprintf(w, "%s%s %s%s%s%s%s\n", tsStr, lvl, fileLine, color, e.Msg, GetAttributes(string(line)), log.Colors.Reset)
}
