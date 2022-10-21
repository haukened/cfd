package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "cfd",
		Usage:    "A lightweight Cloudflare DDNS client, written in Go.",
		Version:  "v0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "David Haukeness",
				Email: "david@hauken.us",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				EnvVars: []string{"CFD_CONFIG"},
				Value:   "/etc/cfd.yml",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Print a bunch of stuff",
				EnvVars: []string{"CFD_DEBUG"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "oneshot",
				Aliases: []string{"now"},
				Usage:   "Performs a one-time update, regardless of history",
				EnvVars: []string{"CFD_ONESHOT"},
				Value:   false,
			},
		},
		Action: run,
	}

	// create a means to gracefully terminate
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
		log.Println("caught SIGINT, exiting.")
	}()

	// listen for hard and soft exits
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(1)
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	debug = c.Bool("debug")
	configFile := c.String("config")
	oneshot := c.Bool("oneshot")
	conf, err := ReadConfig(configFile)
	if err != nil {
		return err
	}
	// perform first update
	err = UpdateCloudflareIP(c.Context, conf)
	if err != nil {
		return err
	}
	// if this isn't a oneshot, start the update loop
	if !oneshot {
		ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)
		for {
			select {
			case <-c.Done():
				return nil
			case <-ticker.C:
				err := UpdateCloudflareIP(c.Context, conf)
				if err != nil {
					log.Printf("unable to update cloudflare IP: %v", err)
				}
			}
		}
	}
	return nil
}
