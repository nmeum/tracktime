package main

import (
	"github.com/nmeum/tracktime/parser"

	"flag"
	"fmt"
	"log"
	"os"
	"time"
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
)

var dateLayout string

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

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
	var maxdurlen int

	workmap := make(map[string]time.Duration)
	for _, entry := range entries {
		key := intervalString(entry.Date)
		if _, ok := workmap[key]; !ok {
			keys = append(keys, key)
		}
		workmap[key] += entry.Duration

		// Date should always have the same width. Only duration
		// requires padding based on the maximum duration length.
		maxdurlen = max(maxdurlen, len(fmt.Sprintf("%v", workmap[key])))
	}

	var delta, goalHours time.Duration
	goalHours = time.Duration(*goal) * time.Hour

	// Output in same order as specified in input file
	for _, key := range keys {
		hours := workmap[key]
		delta += (hours - goalHours)

		// Output should always be aligned at the pipe character.
		fmt.Printf("%v %*v | %v\n", key, maxdurlen, hours, durationString(delta))
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify a file to parse\n")
		os.Exit(1)
	}

	fp := flag.Arg(0)
	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dateLayout = parser.DefaultTimeFormat()
	p := parser.NewParser(dateLayout)

	entries, err := p.ParseEntries(fp, file)
	if err != nil {
		log.Fatal(err)
	}

	handleEntries(entries)
}
