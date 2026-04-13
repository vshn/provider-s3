package bucket

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidator_ValidateCreate_RequireProviderConfig(t *testing.T) {
	tests := map[string]struct {
		providerName  string
		expectedError string
	}{
		"GivenProviderName_ThenExpectNoError": {
			providerName: "provider-config",
		},
		"GivenNoProviderName_ThenExpectError": {
			providerName:  "",
			expectedError: `.spec.providerConfigRef.name is required`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: tc.providerName,
						},
					},
					ForProvider: s3v1.BucketParameters{BucketName: "bucket"},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateCreate(context.TODO(), bucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_PreventBucketNameChange(t *testing.T) {
	tests := map[string]struct {
		newBucketName string
		oldBucketName string
		expectedError string
	}{
		"GivenNoNameInStatus_WhenNoNameInSpec_ThenExpectNil": {
			oldBucketName: "",
			newBucketName: "",
		},
		"GivenNoNameInStatus_WhenNameInSpec_ThenExpectNil": {
			oldBucketName: "",
			newBucketName: "my-bucket",
		},
		"GivenNameInStatus_WhenNameInSpecSame_ThenExpectNil": {
			oldBucketName: "my-bucket",
			newBucketName: "my-bucket",
		},
		"GivenNameInStatus_WhenNameInSpecEmpty_ThenExpectNil": {
			oldBucketName: "bucket",
			newBucketName: "", // defaults to metadata.name
		},
		"GivenNameInStatus_WhenNameInSpecDifferent_ThenExpectError": {
			oldBucketName: "my-bucket",
			newBucketName: "different",
			expectedError: `spec.forProvider.bucketName: Invalid value: "different": Changing the bucket name is not allowed after creation`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{BucketName: tc.oldBucketName},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "provider-config",
						},
					},
				},
				Status: s3v1.BucketStatus{AtProvider: s3v1.BucketProviderStatus{BucketName: tc.oldBucketName}},
			}
			newBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{BucketName: tc.newBucketName},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "provider-config",
						},
					},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_RequireProviderConfig(t *testing.T) {
	tests := map[string]struct {
		providerConfigName string
		expectedError      string
	}{
		"GivenProviderConfigRefWithName_ThenExpectNoError": {
			providerConfigName: "provider-config",
		},
		"GivenProviderConfigEmptyRef_ThenExpectError": {
			providerConfigName: "",
			expectedError:      `spec.providerConfigRef.name: Invalid value: "null": Provider config is required`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: tc.providerConfigName,
						},
					},
					ForProvider: s3v1.BucketParameters{BucketName: "bucket"},
				},
				Status: s3v1.BucketStatus{AtProvider: s3v1.BucketProviderStatus{BucketName: "bucket"}},
			}
			newBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: tc.providerConfigName,
						},
					},
					ForProvider: s3v1.BucketParameters{BucketName: "bucket"},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_PreventZoneChange(t *testing.T) {
	tests := map[string]struct {
		newZone       string
		oldZone       string
		expectedError string
	}{
		"GivenZoneUnchanged_ThenExpectNil": {
			oldZone: "zone",
			newZone: "zone",
		},
		"GivenZoneChanged_ThenExpectError": {
			oldZone:       "zone",
			newZone:       "different",
			expectedError: `spec.forProvider.region: Invalid value: "different": Changing the region is not allowed after creation`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{Region: tc.oldZone},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "provider-config",
						},
					},
				},
				Status: s3v1.BucketStatus{AtProvider: s3v1.BucketProviderStatus{BucketName: "bucket"}},
			}
			newBucket := &s3v1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: s3v1.BucketSpec{
					ForProvider: s3v1.BucketParameters{Region: tc.newZone},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "provider-config",
						},
					},
				},
				Status: s3v1.BucketStatus{AtProvider: s3v1.BucketProviderStatus{BucketName: "bucket"}},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
