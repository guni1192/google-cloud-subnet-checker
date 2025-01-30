package ip

import (
	"fmt"
	"net"

	"github.com/guni1192/google-cloud-subnet-checker/internal/gcloud"
)

// CheckCIDROverlap checks if the desired CIDR overlaps with any of the existing subnets.
func CheckCIDROverlap(subnets []gcloud.Subnetwork, desiredCIDR string) error {
	_, desiredCIDRBlock, err := net.ParseCIDR(desiredCIDR)
	if err != nil {
		return fmt.Errorf("failed to parse desired CIDR: %v", err)
	}

	for _, subnet := range subnets {

		for _, ipRange := range subnet.IPRanges {
			_, subnetCIDR, err := net.ParseCIDR(ipRange.IPv4Range)
			if err != nil {
				return fmt.Errorf("failed to parse subnet CIDR: %v", err)
			}

			if cidrOverlap(subnetCIDR, desiredCIDRBlock) {
				return fmt.Errorf("CIDR %s overlaps with existing subnet: %v", desiredCIDR, subnet)
			}
		}
	}

	return nil
}

// cidrOverlap checks if two CIDR blocks overlap.
func cidrOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}
