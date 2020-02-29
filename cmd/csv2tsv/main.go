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

func scanAndPrint(in *os.File) error {
	reader := csv.NewReader(in)
	reader.LazyQuotes = true
	reader.Comma = ','

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
	if c.NArg() < 1 {
		return scanAndPrint(os.Stdin)
	}

	filepath := c.Args().Get(0)
	in, err := os.Open(filepath)
	defer in.Close()

	if err != nil {
		return err
	}

	return scanAndPrint(in)
}

func main() {
	app := cli.NewApp()
	app.Name = "csv2tsv"
	app.Usage = "convert csv to tsv"
	app.Version = "1.0.0"
	app.ArgsUsage = "[csv_file]?"

	app.Action = func(c *cli.Context) error {
		return run(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
