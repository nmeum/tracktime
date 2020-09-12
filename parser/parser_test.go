package parser

import (
	"testing"
	"time"
)

const testLayout = "02.01.2006"

func getTime(t *testing.T, input string) time.Time {
	ti, err := time.Parse(testLayout, input)
	if err != nil {
		t.Fatal(err)
	}

	return ti
}

func getDur(t *testing.T, input string) time.Duration {
	d, err := time.ParseDuration(input)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func TestParseEntryPositive(t *testing.T) {
	type TestCase struct {
		Input  string
		Result *Entry
	}

	tests := []TestCase{
		{
			"24.12.2010	1000	1500	Foobar",
			&Entry{getTime(t, "24.12.2010"), getDur(t, "5h"), "Foobar"},
		},
		{
			"01.05.1992	1012	1013	foo.",
			&Entry{getTime(t, "01.05.1992"), getDur(t, "1m"), "foo."},
		},
		{
			"13.12.2023	2112	2342	bla",
			&Entry{getTime(t, "13.12.2023"), getDur(t, "2h30m"), "bla"},
		},
		{
			"02.06.2042	1259	2359	foo bar baz",
			&Entry{getTime(t, "02.06.2042"), getDur(t, "11h"), "foo bar baz"},
		},
	}

	parser := NewParser(testLayout)
	for _, test := range tests {
		entry, err := parser.parseEntry(test.Input)
		if err != nil {
			t.Fatal(err)
		}

		if *entry != *test.Result {
			t.Fatalf("Expected %v - got %v", test.Result, entry)
		}
	}
}
