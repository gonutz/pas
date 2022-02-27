package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/akm/delparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s UNIT_FILE.pas\n", os.Args[0])
		os.Exit(1)
	}
	src := os.Args[1]

	f, err := os.Open(src)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	unit, err := delparser.ParseString(string(b))
	if err != nil {
		fmt.Printf("Error parsing file: %+v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(unit, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling unit: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", out)
}
