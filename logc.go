package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"fortio.org/cli"
	"fortio.org/log"
)

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
	BrightRed = "\033[91m]"
)

// Takes json string level and returns the 1 letter short and color to use for it.
func LevelToColor(levelStr string) (string, string) {
	switch levelStr {
	case "dbug":
		return "D", Gray
	case "trace":
		return "V", Cyan
	case "info":
		return "I", Green
	case "warn":
		return "W", Yellow
	case "err":
		return "E", Red
	case "crit":
		return "C", Purple
	case "fatal":
		return "F", BrightRed
	default:
		log.Critf("Bug/Unknown level %q", levelStr)
	}
	return "?", ""
}

func main() {
	noColorFlag := flag.Bool("no-color", false, "Do not colorize output")
	cli.Main()
	// read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	prevDate := time.UnixMilli(0)
	reset := Reset
	noColor := *noColorFlag
	if noColor {
		reset = ""
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		// json deserialize
		e := log.JSONEntry{}
		err := json.Unmarshal(line, &e)
		if err != nil {
			log.LogVf("Error unmarshalling %q: %v", string(line), err)
			fmt.Printf("! %s\n", string(line))
			continue
		}
		ts := time.UnixMicro(e.TS)
		// Each time the day changes we print a header
		if ts.YearDay() != prevDate.YearDay() {
			fmt.Printf("#### %s ####\n", ts.Format("2006-01-02"))
			prevDate = ts
		}
		tsStr := ts.Format("15:04:05.000000")
		// uppercase level
		lvl, color := LevelToColor(e.Level)
		if noColor {
			color = ""
		}
		fmt.Printf("%s%s %s %s:%d> %s%s\n", color, tsStr, lvl, e.File, e.Line, e.Msg, reset)

	}
}
