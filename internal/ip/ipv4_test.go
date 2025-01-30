package ip

import (
	"testing"
)

func TestCheckCIDROverlap(t *testing.T) {
	tests := []struct {
		name          string
		existingCIDRs []string
		desiredCIDR   string
		expectError   bool
	}{
		{
			name:          "No overlap",
			existingCIDRs: []string{"192.168.1.0/24", "10.0.0.0/8"},
			desiredCIDR:   "172.16.0.0/16",
			expectError:   false,
		},
		{
			name:          "Overlap with existing CIDR",
			existingCIDRs: []string{"192.168.1.0/24", "10.0.0.0/8"},
			desiredCIDR:   "192.168.1.128/25",
			expectError:   true,
		},
		{
			name:          "Exact match overlap",
			existingCIDRs: []string{"192.168.1.0/24"},
			desiredCIDR:   "192.168.1.0/24",
			expectError:   true,
		},
		{
			name:          "Overlap with larger CIDR",
			existingCIDRs: []string{"192.168.0.0/16"},
			desiredCIDR:   "192.168.1.0/24",
			expectError:   true,
		},
		{
			name:          "No overlap with adjacent CIDR",
			existingCIDRs: []string{"192.168.1.0/24"},
			desiredCIDR:   "192.168.2.0/24",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCIDROverlap(tt.existingCIDRs, tt.desiredCIDR)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckCIDROverlap() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
