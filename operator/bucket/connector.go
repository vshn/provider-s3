package bucket

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	providerv1 "github.com/vshn/provider-s3/apis/provider/v1"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ managed.ExternalConnector = &connector{}
var _ managed.ExternalClient = &bucketClient{}

const (
	lockAnnotation = s3v1.Group + "/lock"
	KeyID          = "AWS_ACCESS_KEY_ID"
	SecretID       = "AWS_SECRET_ACCESS_KEY"
)

var (
	errNotBucket = fmt.Errorf("managed resource is not a bucket")
)

type connector struct {
	kube     client.Client
	recorder event.Recorder
	usage    *resource.LegacyProviderConfigUsageTracker
}

type bucketClient struct {
	mc       *minio.Client
	recorder event.Recorder
	KeyID    string
	SecretID string
}

// Connect implements managed.ExternalConnector.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("connecting resource")

	bucket, ok := mg.(*s3v1.Bucket)
	if !ok {
		return nil, errNotBucket
	}

	if err := c.usage.Track(ctx, bucket); err != nil {
		return nil, err
	}

	config, err := c.getProviderConfig(ctx, bucket)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	key := client.ObjectKey{Name: config.Spec.Credentials.APISecretRef.Name, Namespace: config.Spec.Credentials.APISecretRef.Namespace}
	err = c.kube.Get(ctx, key, secret)
	if err != nil {
		return nil, err
	}

	mc, err := c.createS3Client(ctx, string(secret.Data[KeyID]), string(secret.Data[SecretID]), config)
	if err != nil {
		return nil, err
	}

	bc := &bucketClient{
		mc:       mc,
		recorder: c.recorder,
		KeyID:    string(secret.Data[KeyID]),
		SecretID: string(secret.Data[SecretID]),
	}

	parsed, err := url.Parse(config.Spec.Endpoint)
	if err != nil {
		return nil, err
	}
	bucket.Status.Endpoint = parsed.Host
	bucket.Status.EndpointURL = parsed.String()

	return bc, nil
}

func (c *connector) createS3Client(ctx context.Context, keyID, secretID string, config *providerv1.ProviderConfig) (*minio.Client, error) {

	if keyID == "" || secretID == "" {
		return nil, fmt.Errorf("credentials missing, please check the keys")
	}

	parsed, err := url.Parse(config.Spec.Endpoint)
	if err != nil {
		return nil, err
	}

	host := parsed.Host
	if parsed.Host == "" {
		host = parsed.Path // if no scheme is given, it's parsed as a path -.-
	}

	return minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(keyID, secretID, ""),
		Secure: isTLSEnabled(parsed),
	})

}

func (c *connector) getProviderConfig(ctx context.Context, bucket *s3v1.Bucket) (*providerv1.ProviderConfig, error) {
	ref := bucket.GetProviderConfigReference()
	if ref == nil {
		return nil, fmt.Errorf("no provider config reference set")
	}
	config := &providerv1.ProviderConfig{}
	err := c.kube.Get(ctx, client.ObjectKey{Name: ref.Name}, config)
	return config, err
}

// isTLSEnabled returns false if the scheme is explicitly set to `http` or `HTTP`
func isTLSEnabled(u *url.URL) bool {
	// if no scheme given, we'll assume tls
	if u.Scheme == "" {
		return true
	}
	return strings.EqualFold(u.Scheme, "https")
}
