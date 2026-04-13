package bucket

import (
	"context"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (b *bucketClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("updating resource")

	bucket, ok := mg.(*s3v1.Bucket)
	if !ok {
		return managed.ExternalUpdate{}, errNotBucket
	}

	policy := ""
	if bucket.Spec.ForProvider.Policy != nil {
		policy = *bucket.Spec.ForProvider.Policy
	}

	// passing an empty string will make minio client remove the bucket policy
	err := b.mc.SetBucketPolicy(ctx, bucket.GetBucketName(), policy)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{}, nil
}
