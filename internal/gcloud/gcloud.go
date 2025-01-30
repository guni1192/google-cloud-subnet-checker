package gcloud

import (
	"context"
	"fmt"
	"log/slog"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// Client wraps the Google Cloud Compute client.
type Client struct {
	subnetClient *compute.SubnetworksClient
}

type Subnetwork struct {
	Name     string
	Region   string
	IPRanges []IPRange
}

type IPRange struct {
	RangeName string
	IPv4Range string
}

// NewClient creates a new Google Cloud Compute client.
func NewClient(ctx context.Context) (*Client, error) {
	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute subnetworks client: %v", err)
	}
	return &Client{subnetClient: subnetClient}, nil
}

// Close closes the Google Cloud Compute client.
func (c *Client) Close() error {
	c.subnetClient.Close()
	return nil
}

// ListSubnetworks retrieves the list of subnetworks for a given project and network.
func (c *Client) ListSubnetworks(ctx context.Context, projectID, region string) ([]Subnetwork, error) {

	req := &computepb.ListSubnetworksRequest{
		Project: projectID,
		Region:  region,
	}

	var subnets []Subnetwork
	it := c.subnetClient.List(ctx, req)
	for {
		subnet, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list subnets: %v", err)
		}

		var ipRanges []IPRange
		ipRanges = append(ipRanges, IPRange{
			RangeName: "primary",
			IPv4Range: *subnet.IpCidrRange,
		})
		slog.Debug("subnet primary ip range", "name", *subnet.Name, "ipRange", *subnet.IpCidrRange)
		for _, ipRange := range subnet.SecondaryIpRanges {
			slog.Debug("subnet secondary ip range", "name", *subnet.Name, "rangeName", *ipRange.RangeName, "ipRange", *ipRange.IpCidrRange)
			ipRanges = append(ipRanges, IPRange{
				RangeName: *ipRange.RangeName,
				IPv4Range: *ipRange.IpCidrRange,
			})
		}

		subnets = append(subnets, Subnetwork{
			Name:     *subnet.Name,
			Region:   *subnet.Region,
			IPRanges: ipRanges,
		})
	}

	return subnets, nil
}
