package bucket

import (
	"strings"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	providerv1 "github.com/vshn/provider-s3/apis/provider/v1"
	s3v1 "github.com/vshn/provider-s3/apis/s3/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetupController adds a controller that reconciles managed resources.
func SetupController(mgr ctrl.Manager) error {
	name := strings.ToLower(s3v1.BucketGroupKind)
	recorder := event.NewAPIRecorder(mgr.GetEventRecorderFor(name)) //nolint:staticcheck // crossplane-runtime v2 NewAPIRecorder doesn't support the new events.EventRecorder yet

	return SetupControllerWithConnector(mgr, name, recorder, &connector{
		kube:     mgr.GetClient(),
		recorder: recorder,
		usage:    resource.NewLegacyProviderConfigUsageTracker(mgr.GetClient(), &providerv1.ProviderConfigUsage{}),
	}, 0*time.Second)
}

func SetupControllerWithConnector(mgr ctrl.Manager, name string, recorder event.Recorder, c managed.ExternalConnector, creationGracePeriod time.Duration) error {
	r := createReconciler(mgr, name, recorder, c, creationGracePeriod)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&s3v1.Bucket{}).
		Complete(r)
}

func createReconciler(mgr ctrl.Manager, name string, recorder event.Recorder, c managed.ExternalConnector, creationGracePeriod time.Duration) *managed.Reconciler {
	return managed.NewReconciler(mgr,
		resource.ManagedKind(s3v1.BucketGroupVersionKind),
		managed.WithExternalConnector(c),
		managed.WithLogger(logging.NewLogrLogger(mgr.GetLogger().WithValues("controller", name))),
		managed.WithRecorder(recorder),
		managed.WithPollInterval(1*time.Minute),
		managed.WithCreationGracePeriod(creationGracePeriod))
}

// SetupWebhook adds a webhook for managed resources.
func SetupWebhook(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &s3v1.Bucket{}).
		WithValidator(&Validator{
			log: mgr.GetLogger().WithName("webhook").WithName(strings.ToLower(s3v1.BucketKind)),
		}).
		Complete()
}
