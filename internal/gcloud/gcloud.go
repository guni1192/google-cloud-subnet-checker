package gcloud

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// Client wraps the Google Cloud Compute client.
type Client struct {
	subnetClient *compute.SubnetworksClient
}

// NewClient creates a new Google Cloud Compute client.
func NewClient(ctx context.Context) (*Client, error) {
	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %v", err)
	}
	return &Client{subnetClient: subnetClient}, nil
}

// Close closes the Google Cloud Compute client.
func (c *Client) Close() error {
	return c.subnetClient.Close()
}

// ListSubnetworks retrieves the list of subnetworks for a given project and network.
func (c *Client) ListSubnetworks(ctx context.Context, projectID, network string) ([]string, error) {
	filter := fmt.Sprintf("network:projects/%s/global/networks/%s", projectID, network)
	req := &computepb.ListSubnetworksRequest{
		Project: projectID,
		Filter:  &filter,
	}

	var existingCIDRs []string
	it := c.subnetClient.List(ctx, req)
	for {
		subnet, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list subnets: %v", err)
		}

		existingCIDRs = append(existingCIDRs, *subnet.IpCidrRange)
	}

	return existingCIDRs, nil
}
