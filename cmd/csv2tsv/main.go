package main

import (
	"encoding/csv"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"strings"
)

func scanAndPrint(in *os.File, force bool) error {
	reader := csv.NewReader(in)
	reader.LazyQuotes = true
	reader.Comma = ','
	if force {
		reader.FieldsPerRecord = -1
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Println(strings.Join(record, "\t"))

		if err != nil {
			return err
		}
	}

	return nil
}

func run(c *cli.Context) error {
	force := c.Bool("force")

	if c.NArg() < 1 {
		return scanAndPrint(os.Stdin, force)
	}

	filepath := c.Args().Get(0)
	in, err := os.Open(filepath)
	defer in.Close()

	if err != nil {
		return err
	}
	return scanAndPrint(in, force)
}

func main() {
	app := cli.NewApp()
	app.Name = "csv2tsv"
	app.Usage = "convert csv to tsv"
	app.Version = "1.0.0"
	app.ArgsUsage = "[csv_file]?"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "force",
			Aliases:  []string{"f"},
			Usage:    "force skip checking each column size",
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
