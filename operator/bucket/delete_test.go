package bucket

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
)

func TestDelete_NotABucket(t *testing.T) {
	b := bucketClient{}
	ctx := logr.NewContext(context.Background(), logr.Discard())
	_, err := b.Delete(ctx, &fake.Managed{})
	assert.EqualError(t, err, errNotBucket.Error())
}

func TestHasDeleteAllPolicy(t *testing.T) {
	tests := map[string]struct {
		policy   s3v1.BucketDeletionPolicy
		expected bool
	}{
		"DeleteAll": {
			policy:   s3v1.DeleteAll,
			expected: true,
		},
		"DeleteIfEmpty": {
			policy:   s3v1.DeleteIfEmpty,
			expected: false,
		},
		"Empty": {
			policy:   "",
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bucket := &s3v1.Bucket{
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{
						BucketDeletionPolicy: tc.policy,
					},
				},
			}
			assert.Equal(t, tc.expected, hasDeleteAllPolicy(bucket))
		})
	}
}
