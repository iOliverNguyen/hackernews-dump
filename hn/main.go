package hn

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/iOliverNguyen/hackernews/log"
)

var (
	ItemsPerFile = 1000
	FilesPerDir  = 1000

	Concurrent     = 30
	SaveCheckpoint = 100
	DataDir        = "data"
	UpdateBack     = 0

	logger      = log.NewStdTextLogger(os.Stderr, 0)
	client      = NewHNClient()
	projectRoot string
	flagDebug   bool
)

func setupDebug() {
	ItemsPerFile = 10
	Concurrent = 1
	SaveCheckpoint = 5
	DataDir = "data-debug"
}

func setup() {
	projectRoot = os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		panic("no PROJECT_ROOT")
	}

	must(0, os.MkdirAll(projectRoot+"/"+DataDir, 0755))
}

func Main() {
	app := cli.App{
		Name:        "hn",
		Description: "HackerNews tool",
		Action: func(ctx *cli.Context) error {
			fmt.Println("HackerNews tool")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "sync",
				Usage:  "sync data from hacker news into files",
				Action: cmdSyncFiles,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "debug",
						Usage:       "enable debug mode, set concurrent to 2, items_per_files to 10, and save checkpoint to 5",
						Value:       false,
						Destination: &flagDebug,
					},
					&cli.IntFlag{
						Name:        "update-back",
						Usage:       "update files that are newer than this number of days, for refetching updated scores",
						Value:       0,
						Destination: &UpdateBack,
					},
					&cli.IntFlag{
						Name:        "concurrent",
						Value:       30,
						Destination: &Concurrent,
					},
				},
			},
			{
				Name:   "load",
				Usage:  "load data from files into memory",
				Action: cmdLoadMem,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "concurrent",
						Value:       30,
						Destination: &Concurrent,
					},
				},
			},
		},
		Before: func(context *cli.Context) error {
			setup()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Log(log.ErrorLevel, fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
