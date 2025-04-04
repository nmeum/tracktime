package parser

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

const lineFormat = "^\\+?(..*)	([0-9][0-9][0-9][0-9])	([0-9][0-9][0-9][0-9])	(..*)$"

const (
	layoutEnv     = "TRACKTIME_FORMAT"
	defaultLayout = "02.01.2006"
)

type Entry struct {
	Date        time.Time
	Duration    time.Duration
	Description string
	BonusWork   bool
}

type Parser struct {
	validLine *regexp.Regexp
	layout    string
	lineNum   uint
}

// https://en.wikipedia.org/wiki/24-hour_clock#Military_time
func (p *Parser) militaryTime(time int) (int, error) {
	hours := time / 100
	minutes := time % 100

	if hours > 24 {
		return 0, errors.New("invalid hour in duration")
	} else if minutes >= 60 {
		return 0, errors.New("invalid minute in duration")
	}

	return (hours * 60) + minutes, nil
}

func (p *Parser) getFields(line string) (bool, string, int, int, string, error) {
	matches := p.validLine.FindStringSubmatch(line)
	if matches == nil {
		return false, "", 0, 0, "", errors.New("line does not match format")
	}

	durStart, err := strconv.Atoi(matches[2])
	if err != nil {
		return false, "", 0, 0, "", err
	}
	durEnd, err := strconv.Atoi(matches[3])
	if err != nil {
		return false, "", 0, 0, "", err
	}

	// If the line starts with a “+” character, then count this as
	// bonus working time (e.g. working on the weekend).
	bonusEntry := matches[0][0] == '+'

	return bonusEntry, matches[1], durStart, durEnd, matches[4], nil
}

func (p *Parser) parseEntry(line string) (*Entry, error) {
	bonus, date, durStart, durEnd, desc, err := p.getFields(line)
	if err != nil {
		return nil, err
	}

	etime, err := time.Parse(p.layout, date)
	if err != nil {
		return nil, err
	}

	if durStart >= durEnd {
		return nil, errors.New("invalid duration")
	}
	start, err := p.militaryTime(durStart)
	if err != nil {
		return nil, err
	}
	end, err := p.militaryTime(durEnd)
	if err != nil {
		return nil, err
	}

	// Add start duration to entry date, allows reconstructing the
	// absolute start time (and end time) of the given entry.
	etime = etime.Add(time.Duration(start) * time.Minute)

	duration := time.Duration(end-start) * time.Minute
	return &Entry{etime, duration, desc, bonus}, nil
}

func (p *Parser) ParseEntries(fn string, r io.Reader) ([]*Entry, error) {
	var entries []*Entry

	// Reset line number information
	p.lineNum = 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p.lineNum++
		line := scanner.Text()

		entry, err := p.parseEntry(line)
		if err != nil {
			return entries, ParserError{fn, p.lineNum, err.Error()}
		}

		entries = append(entries, entry)
	}

	err := scanner.Err()
	if err != nil {
		return entries, err
	}

	return entries, nil
}

func NewParser(layout string) *Parser {
	validLine := regexp.MustCompile(lineFormat)
	return &Parser{validLine, layout, 0}
}

func DefaultTimeFormat() string {
	dateLayout := os.Getenv(layoutEnv)
	if dateLayout == "" {
		dateLayout = defaultLayout
	}

	return dateLayout
}
