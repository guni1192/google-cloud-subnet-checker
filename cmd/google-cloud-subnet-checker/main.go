package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/guni1192/google-cloud-subnet-checker/internal/gcloud"
	"github.com/guni1192/google-cloud-subnet-checker/internal/ip"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "google-cloud-subnet-checker",
		Usage: "Check for CIDR overlap in a Google Cloud network before creating a new subnet",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Usage:    "Google Cloud Project ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "region",
				Usage:    "Region name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cidr",
				Usage:    "Desired CIDR for the new subnet",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
		},
		Action: func(c *cli.Context) error {
			projectID := c.String("project")
			region := c.String("region")
			desiredCIDR := c.String("cidr")
			debug := c.Bool("debug")

			// Set up slog with appropriate log level
			var logger *slog.Logger
			if debug {
				logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			} else {
				logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			}
			slog.SetDefault(logger)

			ctx := context.Background()
			client, err := gcloud.NewClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			existingCIDRs, err := client.ListSubnetworks(ctx, projectID, region)
			if err != nil {
				return err
			}

			if err := ip.CheckCIDROverlap(existingCIDRs, desiredCIDR); err != nil {
				return err
			}

			fmt.Println("No CIDR overlap detected. Safe to create new subnet.")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
	}
}
