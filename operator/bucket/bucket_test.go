package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetBucketName(t *testing.T) {
	tests := map[string]struct {
		metadataName string
		specName     string
		expected     string
	}{
		"UsesSpecBucketNameIfSet": {
			metadataName: "resource-name",
			specName:     "custom-bucket",
			expected:     "custom-bucket",
		},
		"FallsBackToMetadataName": {
			metadataName: "resource-name",
			specName:     "",
			expected:     "resource-name",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: tc.metadataName},
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{
						BucketName: tc.specName,
					},
				},
			}
			assert.Equal(t, tc.expected, bucket.GetBucketName())
		})
	}
}
