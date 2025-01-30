package main

import (
	"context"
	"fmt"
	"log"
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
				Name:     "project-id",
				Usage:    "Google Cloud Project ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "network",
				Usage:    "Network name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "desired-cidr",
				Usage:    "Desired CIDR for the new subnet",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			projectID := c.String("project-id")
			network := c.String("network")
			desiredCIDR := c.String("desired-cidr")

			ctx := context.Background()
			client, err := gcloud.NewClient(ctx)
			if err != nil {
				return err
			}
			defer client.Close()

			existingCIDRs, err := client.ListSubnetworks(ctx, projectID, network)
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
		log.Fatal(err)
	}
}
