package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	defaultLayout = "02.01.2006" // TODO: Make configurable through env
	lineFormat    = "%s	%04d	%04d	%s"
)

var (
	hours    = flag.Int("h", 8, "hours per interval")
	interval = flag.String("i", "d", "interval for working hours")
)

type ParserError struct {
	Line uint
	Msg  string
}

func (p ParserError) Error() string {
	return fmt.Sprintf("%s:%d %s", "stdin", p.Line, p.Msg)
}

type Entry struct {
	Date        time.Time
	Duration    time.Duration
	Description string
}

// https://en.wikipedia.org/wiki/24-hour_clock#Military_time
func militaryTime(time uint) (uint, error) {
	hours := time / 100
	minutes := time % 100

	if hours > 24 {
		return 0, fmt.Errorf("invalid hour")
	} else if minutes >= 60 {
		return 0, fmt.Errorf("invalid minute")
	}

	return (hours * 60) + minutes, nil
}

func parseEntry(line string) (entry Entry, err error) {
	var date, desc string
	var durStart, durEnd uint

	_, err = fmt.Sscanf(line, lineFormat, &date, &durStart, &durEnd, &desc)
	if err != nil {
		return Entry{}, err
	}
	entry.Description = desc

	entry.Date, err = time.Parse(defaultLayout, date)
	if err != nil {
		return Entry{}, err
	}

	if durStart >= durEnd {
		return Entry{}, ParserError{23, "invalid duration"}
	}
	start, err := militaryTime(durStart)
	if err != nil {
		return Entry{}, ParserError{23, "invalid start duration: " + err.Error()}
	}
	end, err := militaryTime(durEnd)
	if err != nil {
		return Entry{}, ParserError{23, "invalid end duration: " + err.Error()}
	}

	entry.Duration = time.Duration(end-start) * time.Minute
	return entry, nil
}

func parseEntries(r io.Reader) (entries []Entry, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseEntry(line)
		if err != nil {
			return entries, err
		}

		entries = append(entries, entry)
	}

	err = scanner.Err()
	if err != nil {
		return entries, err
	}

	return entries, nil
}

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
