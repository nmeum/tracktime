package main

import (
	"github.com/nmeum/tracktime/parser"

	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Table [][]string

const outputFormat = time.RFC822

func printTabular(w io.Writer, rows Table) {
	numCols := len(rows[0])

	alignment := make([]string, numCols)
	for i := 0; i < numCols; i++ {
		alignment[i] = "r"
	}
	alignment[0] = "l"

	fmt.Fprintf(w, "\\begin{tabular}{%s}\n", strings.Join(alignment, "|"))
	for _, row := range rows {
		fmt.Fprintf(w, "\t")
		for n, col := range row {
			delim := "&"
			if n == numCols-1 {
				delim = "\\\\"
			}

			fmt.Fprintf(w, "%v %v ", col, delim)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "\\end{tabular}\n")
}

func main() {
	p := parser.NewParser(parser.DefaultTimeFormat())
	entries, err := p.ParseEntries("stdin", os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var rows Table
	for _, entry := range entries {
		rows = append(rows, []string{
			entry.Date.Format(outputFormat),
			entry.Duration.String(),
			entry.Description,
		})
	}

	if len(rows) == 0 {
		log.Fatal("no entries in given file")
	}
	printTabular(os.Stdout, rows)
}
