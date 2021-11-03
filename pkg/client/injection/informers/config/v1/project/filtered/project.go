// Code generated by injection-gen. DO NOT EDIT.

package filtered

import (
	context "context"

	versioned "github.com/openshift-knative/serverless-operator/pkg/client/clientset/versioned"
	v1 "github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1"
	client "github.com/openshift-knative/serverless-operator/pkg/client/injection/client"
	filtered "github.com/openshift-knative/serverless-operator/pkg/client/injection/informers/factory/filtered"
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
	injection.Default.RegisterFilteredInformers(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct {
	Selector string
}

func withInformer(ctx context.Context) (context.Context, []controller.Informer) {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	infs := []controller.Informer{}
	for _, selector := range labelSelectors {
		f := filtered.Get(ctx, selector)
		inf := f.Config().V1().Projects()
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
		infs = append(infs, inf.Informer())
	}
	return ctx, infs
}

func withDynamicInformer(ctx context.Context) context.Context {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	for _, selector := range labelSelectors {
		inf := &wrapper{client: client.Get(ctx), selector: selector}
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
	}
	return ctx
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context, selector string) v1.ProjectInformer {
	untyped := ctx.Value(Key{Selector: selector})
	if untyped == nil {
		logging.FromContext(ctx).Panicf(
			"Unable to fetch github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1.ProjectInformer with selector %s from context.", selector)
	}
	return untyped.(v1.ProjectInformer)
}

type wrapper struct {
	client versioned.Interface

	selector string
}

var _ v1.ProjectInformer = (*wrapper)(nil)
var _ configv1.ProjectLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiconfigv1.Project{}, 0, nil)
}

func (w *wrapper) Lister() configv1.ProjectLister {
	return w
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiconfigv1.Project, err error) {
	reqs, err := labels.ParseToRequirements(w.selector)
	if err != nil {
		return nil, err
	}
	selector = selector.Add(reqs...)
	lo, err := w.client.ConfigV1().Projects().List(context.TODO(), metav1.ListOptions{
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

func (w *wrapper) Get(name string) (*apiconfigv1.Project, error) {
	// TODO(mattmoor): Check that the fetched object matches the selector.
	return w.client.ConfigV1().Projects().Get(context.TODO(), name, metav1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
