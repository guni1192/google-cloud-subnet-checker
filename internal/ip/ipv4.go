package ip

import (
	"fmt"
	"net"
)

// CheckCIDROverlap checks if the desired CIDR overlaps with any of the existing subnets.
func CheckCIDROverlap(existingCIDRs []string, desiredCIDR string) error {
	_, desiredCIDRBlock, err := net.ParseCIDR(desiredCIDR)
	if err != nil {
		return fmt.Errorf("failed to parse desired CIDR: %v", err)
	}

	for _, existingCIDR := range existingCIDRs {
		_, subnetCIDR, err := net.ParseCIDR(existingCIDR)
		if err != nil {
			return fmt.Errorf("failed to parse subnet CIDR: %v", err)
		}

		if cidrOverlap(subnetCIDR, desiredCIDRBlock) {
			return fmt.Errorf("CIDR %s overlaps with existing subnet %s", desiredCIDR, existingCIDR)
		}
	}

	return nil
}

// cidrOverlap checks if two CIDR blocks overlap.
func cidrOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}
