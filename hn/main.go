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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Destination: &flagDebug,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("HackerNews tool")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "sync",
				Usage:  "sync data from hacker news into files",
				Action: cmdSyncFiles,
			},
		},
		Before: func(context *cli.Context) error {
			if flagDebug {
				setupDebug()
			}
			setup()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Log(log.ErrorLevel, fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
