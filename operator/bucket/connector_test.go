package bucket

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTLSEnabled(t *testing.T) {
	tests := map[string]struct {
		rawURL   string
		expected bool
	}{
		"HTTPS": {
			rawURL:   "https://s3.example.com",
			expected: true,
		},
		"HTTP": {
			rawURL:   "http://s3.example.com",
			expected: false,
		},
		"HTTPS_uppercase": {
			rawURL:   "HTTPS://s3.example.com",
			expected: true,
		},
		"HTTP_uppercase": {
			rawURL:   "HTTP://s3.example.com",
			expected: false,
		},
		"NoScheme": {
			rawURL:   "s3.example.com",
			expected: true,
		},
		"MixedCase_HttpS": {
			rawURL:   "HttpS://s3.example.com",
			expected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tc.rawURL)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expected, isTLSEnabled(u))
		})
	}
}
