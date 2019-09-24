package postory_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupported(t *testing.T) {
	testcases := []struct {
		carrierString string
		expected      bool
	}{
		{"canada_post", true},
		{"dhl_express", true},
		{"fedex", true},
		{"ups", true},
		{"usps", true},
		{"fake_carrier", false},
	}

	for _, tc := range testcases {
		actual := isCarrierSupported(tc.carrierString)
		assert.Equal(t, tc.expected, actual)
	}
}
