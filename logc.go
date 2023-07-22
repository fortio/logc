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

// Takes json string level and returns the 1 letter short and color to use for it.
func LevelToColor(levelStr string) (string, string) {
	level, found := log.JSONStringLevelToLevel[levelStr]
	if !found {
		log.Critf("Bug/Unknown level %q", levelStr)
		return "?", log.Colors.Blue
	}
	return log.LevelToStrA[level][0:1], log.LevelToColor[level]
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
	noColorFlag := flag.Bool("no-color", false, "Do not colorize output")
	cli.ArgsHelp = " < log.json\nor for instance\n\tfortio server 2>&1 | logc\n" +
		"to convert JSON fortio logger lines from stdin to (ansi) colorized text"
	cli.Main()
	// read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	prevDate := time.UnixMilli(0)
	noColor := *noColorFlag
	for scanner.Scan() {
		line := scanner.Bytes()
		ProcessLogLine(os.Stdout, &prevDate, noColor, line)
	}
}

func ProcessLogLine(w io.Writer, prevDate *time.Time, noColor bool, line []byte) {
	// json deserialize
	e := log.JSONEntry{}
	err := json.Unmarshal(line, &e)
	if err != nil {
		log.LogVf("Error unmarshalling %q: %v", string(line), err)
		fmt.Printf("! %s\n", string(line))
		return
	}
	tsStr := ""
	reset := log.Colors.Reset
	darkGrey := log.Colors.DarkGray
	grey := log.Colors.Gray
	// uppercase single letter level + color extraction
	lvl, color := LevelToColor(e.Level)
	if noColor {
		color = ""
		reset = ""
		darkGrey = ""
		grey = ""
	}
	if e.TS != 0 {
		ts := e.Time()
		// Each time the day changes we print a header
		if ts.YearDay() != prevDate.YearDay() {
			fmt.Printf("#### %s ####\n", ts.Format("2006-01-02"))
			*prevDate = ts
		}
		// Use full microseconds resolution unlike the log built in color version which stops at millis.
		tsStr = ts.Format(darkGrey + "15:04:05.000000 ")
	}
	if e.R > 0 {
		tsStr += fmt.Sprintf(grey+"[%d] ", e.R)
	}
	fileLine := ""
	if e.Line != 0 {
		fileLine = fmt.Sprintf("%s:%d> ", e.File, e.Line)
	}
	// Msg can be multi line.
	fmt.Fprintf(w, "%s%s%s %s%s%s%s\n", tsStr, color, lvl, fileLine, e.Msg, GetAttributes(string(line)), reset)
}
