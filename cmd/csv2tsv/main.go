package main

import (
	"encoding/csv"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func scanAndPrint(in *os.File, uncheckFieldSize bool, pretty bool, emptyFillString string) error {
	reader := csv.NewReader(in)
	reader.LazyQuotes = true
	reader.Comma = ','
	if uncheckFieldSize {
		reader.FieldsPerRecord = -1
	}
	blankRegex := regexp.MustCompile("^\\s*$")
	prettyPrefixRegex := regexp.MustCompile("^\\s+(.+)$")
	prettySuffixRegex := regexp.MustCompile("^(.+)\\s+$")

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if len(emptyFillString) > 0 {
			for i := 0; i < len(record); i++ {
				if pretty {
					record[i] = strings.ReplaceAll(record[i], "\n", " ")
					record[i] = prettyPrefixRegex.ReplaceAllString(record[i], "$1")
					record[i] = prettySuffixRegex.ReplaceAllString(record[i], "$1")
				}

				if blankRegex.MatchString(record[i]) {
					record[i] = emptyFillString
				}
			}
		}

		fmt.Println(strings.Join(record, "\t"))

		if err != nil {
			return err
		}
	}

	return nil
}

func run(c *cli.Context) error {
	force := c.Bool("skip")
	emptyFillString := c.String("empty")
	pretty := c.Bool("pretty")

	if c.NArg() < 1 {
		return scanAndPrint(os.Stdin, force, pretty, emptyFillString)
	}

	filepath := c.Args().Get(0)
	in, err := os.Open(filepath)
	defer in.Close()

	if err != nil {
		return err
	}
	return scanAndPrint(in, force, pretty, emptyFillString)
}

func main() {
	app := cli.NewApp()
	app.Name = "csv2tsv"
	app.Usage = "convert csv to tsv"
	app.Version = "1.0.0"
	app.ArgsUsage = "[csv_file]?"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "skip",
			Aliases:  []string{"s"},
			Usage:    "skip checking each column size",
			Value:    false,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "empty",
			Aliases:  []string{"e"},
			Usage:    "fill blank column by specific string. e.g. `-e NaN` replace \" \" to NaN.",
			Value:    "",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "pretty",
			Aliases:  []string{"p"},
			Usage:    "remove prefix/suffix blank string and replace new line to space",
			Value:    false,
			Required: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		return run(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
