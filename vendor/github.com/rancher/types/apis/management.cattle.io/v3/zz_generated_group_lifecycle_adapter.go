package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type GroupLifecycle interface {
	Create(obj *Group) (runtime.Object, error)
	Remove(obj *Group) (runtime.Object, error)
	Updated(obj *Group) (runtime.Object, error)
}

type groupLifecycleAdapter struct {
	lifecycle GroupLifecycle
}

func (w *groupLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *groupLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Group))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGroupLifecycleAdapter(name string, clusterScoped bool, client GroupInterface, l GroupLifecycle) GroupHandlerFunc {
	adapter := &groupLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Group) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
