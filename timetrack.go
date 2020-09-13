package main

import (
	"github.com/nmeum/timetrack/parser"

	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	layoutEnv     = "TIMETRACK_FORMAT"
	defaultLayout = "02.01.2006"
)

const (
	DAY   = 'd'
	WEEK  = 'w'
	MONTH = 'm'
)

var (
	goal     = flag.Int("h", 8, "hours per interval")
	interval = flag.String("i", "d", "interval for working hours")
	seconds  = flag.Bool("s", false, "output duration in seconds")
	total    = flag.Bool("t", false, "output total delta")
)

var dateLayout string

func durationString(duration time.Duration) string {
	if *seconds {
		return fmt.Sprintf("%v", duration.Seconds())
	} else {
		return duration.String()
	}
}

func intervalString(date time.Time) string {
	if *interval == "" {
		fmt.Fprintf(os.Stderr, "invalid interval\n")
		os.Exit(1)
	}

	switch (*interval)[0] {
	case DAY:
		return date.Format(dateLayout)
	case WEEK:
		year, week := date.ISOWeek()
		return fmt.Sprintf("W%v %v", week, year)
	case MONTH:
		year := date.Year()
		return fmt.Sprintf("%s %v", date.Month(), year)
	default:
		fmt.Fprintf(os.Stderr, "unsupported interval: %q\n", *interval)
		os.Exit(2)
	}

	panic("unreachable")
}

func handleEntries(entries []*parser.Entry) {
	var keys []string

	workmap := make(map[string]time.Duration)
	for _, entry := range entries {
		key := intervalString(entry.Date)
		workmap[key] += entry.Duration

		keys = append(keys, key)
	}

	var delta, goalHours time.Duration
	goalHours = time.Duration(*goal) * time.Hour

	// Output in same order as specified in input file
	for _, key := range keys {
		hours := workmap[key]
		delta += (hours - goalHours)

		fmt.Printf("%v\t%v\t| %v\n", key, hours, durationString(delta))
	}

	if *total {
		fmt.Printf("\n---\n\nCurrent overall delta: %v\n", durationString(delta))
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify a file to parse\n")
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	p := parser.NewParser(defaultLayout)
	entries, err := p.ParseEntries(file)
	if err != nil {
		log.Fatal(err)
	}

	dateLayout = os.Getenv(layoutEnv)
	if dateLayout == "" {
		dateLayout = defaultLayout
	}

	handleEntries(entries)
}
