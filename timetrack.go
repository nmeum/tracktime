package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	hours    = flag.Int("h", 8, "hours per interval")
	interval = flag.String("i", "d", "interval for working hours")
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify a file to parse")
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	entries, err := parseEntries(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}
