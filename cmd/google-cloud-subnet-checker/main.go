package main

import (
	"context"
	"fmt"
	"log"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/iterator"

	"github.com/guni1192/google-cloud-subnet-checker/internal/ip"
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
			client, err := compute.NewSubnetworksRESTClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create compute client: %v", err)
			}
			defer client.Close()

			filter := fmt.Sprintf("network:projects/%s/global/networks/%s", projectID, network)
			req := &computepb.ListSubnetworksRequest{
				Project: projectID,
				Filter:  &filter,
			}

			var existingCIDRs []string
			it := client.List(ctx, req)
			for {
				subnet, err := it.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return fmt.Errorf("failed to list subnets: %v", err)
				}

				existingCIDRs = append(existingCIDRs, *subnet.IpCidrRange)
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
