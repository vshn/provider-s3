package bucket

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func TestUpdate_NotABucket(t *testing.T) {
	b := bucketClient{}
	ctx := logr.NewContext(context.Background(), logr.Discard())
	_, err := b.Update(ctx, &fake.Managed{})
	assert.EqualError(t, err, errNotBucket.Error())
}
