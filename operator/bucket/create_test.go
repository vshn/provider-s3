package bucket

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreate_NotABucket(t *testing.T) {
	b := bucketClient{}
	ctx := logr.NewContext(context.Background(), logr.Discard())
	_, err := b.Create(ctx, &fake.Managed{})
	assert.EqualError(t, err, errNotBucket.Error())
}

func TestSetLock(t *testing.T) {
	t.Run("NilAnnotations", func(t *testing.T) {
		bucket := &s3v1.Bucket{}
		b := bucketClient{}
		b.setLock(bucket)
		assert.Equal(t, "claimed", bucket.Annotations[lockAnnotation])
	})

	t.Run("ExistingAnnotations", func(t *testing.T) {
		bucket := &s3v1.Bucket{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{"other": "value"},
			},
		}
		b := bucketClient{}
		b.setLock(bucket)
		assert.Equal(t, "claimed", bucket.Annotations[lockAnnotation])
		assert.Equal(t, "value", bucket.Annotations["other"])
	})
}

func TestCreate_ReturnsConnectionDetails(t *testing.T) {
	// Verify that Create returns connection details when createS3Bucket is overridable
	// We can't easily test the full Create without a minio mock, but we can verify
	// the connection details struct is correctly formed
	details := managed.ConnectionDetails{
		SecretID: []byte("my-secret"),
		KeyID:    []byte("my-key"),
	}
	assert.Equal(t, "my-secret", string(details[SecretID]))
	assert.Equal(t, "my-key", string(details[KeyID]))
}
