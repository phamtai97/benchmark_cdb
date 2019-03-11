package cli

import (
	"benchmark_cockroachdb/configvar"
	"benchmark_cockroachdb/executor"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	DEFAULT_NUMGOROUTINE = 1000
	DEFAULT_NUMMSG       = 1000000
	DEFAULT_TYPE         = 0
)

//Initialize init var
func Initialize() {
	app := cli.NewApp()

	cliConfVar := &configvar.CliConfVar{}

	flags := []cli.Flag{
		cli.IntFlag{
			Name:        "numGoroutine",
			Value:       DEFAULT_NUMGOROUTINE,
			Usage:       "number goroutine run benchmark",
			Destination: &cliConfVar.NumGoroutine,
		},
		cli.Int64Flag{
			Name:        "numMsg",
			Value:       DEFAULT_NUMMSG,
			Usage:       "number msg run benchmark",
			Destination: &cliConfVar.NumMsg,
		},
		cli.IntFlag{
			Name:        "type",
			Value:       DEFAULT_TYPE,
			Usage:       "0: create account. 1: query account",
			Destination: &cliConfVar.Type,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "executor",
			Usage: "benchmark",
			Flags: flags,
			Action: func(c *cli.Context) error {
				executor.Execute(cliConfVar)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
