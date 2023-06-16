package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"fortio.org/cli"
	"fortio.org/log"
)

func main() {
	cli.Main()
	// read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		// json deserialize
		e := log.JSONEntry{}
		err := json.Unmarshal(line, &e)
		if err != nil {
			log.Warnf("Error unmarshalling %q: %v", string(line), err)
		} else {
			fmt.Printf("%+v\n", e)
		}
	}
}
