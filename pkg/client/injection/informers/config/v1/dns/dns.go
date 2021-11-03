// Code generated by injection-gen. DO NOT EDIT.

package dns

import (
	context "context"

	versioned "github.com/openshift-knative/serverless-operator/pkg/client/clientset/versioned"
	v1 "github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1"
	client "github.com/openshift-knative/serverless-operator/pkg/client/injection/client"
	factory "github.com/openshift-knative/serverless-operator/pkg/client/injection/informers/factory"
	configv1 "github.com/openshift-knative/serverless-operator/pkg/client/listers/config/v1"
	apiconfigv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Config().V1().DNSs()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1.DNSInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1.DNSInformer from context.")
	}
	return untyped.(v1.DNSInformer)
}

type wrapper struct {
	client versioned.Interface
}

var _ v1.DNSInformer = (*wrapper)(nil)
var _ configv1.DNSLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiconfigv1.DNS{}, 0, nil)
}

func (w *wrapper) Lister() configv1.DNSLister {
	return w
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiconfigv1.DNS, err error) {
	lo, err := w.client.ConfigV1().DNSs().List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiconfigv1.DNS, error) {
	return w.client.ConfigV1().DNSs().Get(context.TODO(), name, metav1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
