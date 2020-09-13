package parser

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"
)

const lineFormat = "^(..*)	([0-9][0-9][0-9][0-9])	([0-9][0-9][0-9][0-9])	(..*)$"

type Entry struct {
	Date        time.Time
	Duration    time.Duration
	Description string
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
		return 0, errors.New("invalid hour")
	} else if minutes >= 60 {
		return 0, errors.New("invalid minute")
	}

	return (hours * 60) + minutes, nil
}

func (p *Parser) getFields(line string) (string, int, int, string, error) {
	matches := p.validLine.FindStringSubmatch(line)
	if matches == nil {
		return "", 0, 0, "", errors.New("line does not match format")
	}

	durStart, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, 0, "", err
	}
	durEnd, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", 0, 0, "", err
	}

	return matches[1], durStart, durEnd, matches[4], nil
}

func (p *Parser) parseEntry(line string) (*Entry, error) {
	date, durStart, durEnd, desc, err := p.getFields(line)
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

	duration := time.Duration(end-start) * time.Minute
	return &Entry{etime, duration, desc}, nil
}

func (p *Parser) ParseEntries(fn string, r io.Reader) ([]*Entry, error) {
	var entries []*Entry

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
