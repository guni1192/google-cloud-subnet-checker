package ip

import (
	"testing"

	"github.com/guni1192/google-cloud-subnet-checker/internal/gcloud"
)

func TestCheckCIDROverlap(t *testing.T) {
	tests := []struct {
		name        string
		subnets     []gcloud.Subnetwork
		desiredCIDR string
		expectError bool
	}{
		{
			name: "No overlap",
			subnets: []gcloud.Subnetwork{
				{
					Name:     "subnet-1",
					Region:   "us-central1",
					IPRanges: []gcloud.IPRange{{IPv4Range: "192.168.1.0/24"}},
				},
			},
			desiredCIDR: "172.16.0.0/16",
			expectError: false,
		},
		{
			name: "Overlap with existing CIDR",
			subnets: []gcloud.Subnetwork{
				{
					Name:   "subnet-1",
					Region: "us-central1",
					IPRanges: []gcloud.IPRange{
						{IPv4Range: "192.168.1.0/24"},
						{IPv4Range: "10.0.0.0/8"},
					},
				},
			},
			desiredCIDR: "192.168.1.128/25",
			expectError: true,
		},
		{
			name: "Exact match overlap",
			subnets: []gcloud.Subnetwork{
				{
					Name:   "subnet-1",
					Region: "us-central1",
					IPRanges: []gcloud.IPRange{
						{IPv4Range: "192.168.1.0/24"},
					},
				},
			},
			desiredCIDR: "192.168.1.0/24",
			expectError: true,
		},
		{
			name: "Overlap with larger CIDR",
			subnets: []gcloud.Subnetwork{
				{
					Name:   "subnet-1",
					Region: "us-central1",
					IPRanges: []gcloud.IPRange{
						{IPv4Range: "192.168.0.0/16"},
					},
				},
			},
			desiredCIDR: "192.168.1.0/24",
			expectError: true,
		},
		{
			name: "No overlap with adjacent CIDR",
			subnets: []gcloud.Subnetwork{
				{
					Name:   "subnet-1",
					Region: "us-central1",
					IPRanges: []gcloud.IPRange{
						{IPv4Range: "192.168.1.0/24"},
					},
				},
			},
			desiredCIDR: "192.168.2.0/24",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCIDROverlap(tt.subnets, tt.desiredCIDR)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckCIDROverlap() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
