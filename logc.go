package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"fortio.org/cli"
	"fortio.org/log"
)

func main() {
	cli.Main()
	// read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	prevDate := time.UnixMilli(0)
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
		lvl := strings.ToUpper(e.Level[:1])
		fmt.Printf("%s %s %s:%d> %+v\n", tsStr, lvl, e.File, e.Line, e.Msg)

	}
}
