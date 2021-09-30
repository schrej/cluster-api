package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// ObjectReference references a Kubernetes API object by name and namespace.
type ObjectReference struct {
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
}

func (r *ObjectReference) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(r.APIVersion, r.Kind)
}

func (r *ObjectReference) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	r.APIVersion, r.Kind = gvk.ToAPIVersionAndKind()
}

func (r ObjectReference) LocalRef() *LocalObjectReference {
	return &LocalObjectReference{
		Kind:       r.Kind,
		Name:       r.Name,
		APIVersion: r.APIVersion,
	}
}

func (r *ObjectReference) ToCoreRef() *corev1.ObjectReference {
	return &corev1.ObjectReference{
		Kind:       r.Kind,
		APIVersion: r.APIVersion,
		Namespace:  r.Namespace,
		Name:       r.Name,
	}
}

func ObjectReferenceFromCore(ref corev1.ObjectReference) *ObjectReference {
	return &ObjectReference{
		Kind:       ref.Kind,
		Namespace:  ref.Namespace,
		Name:       ref.Name,
		APIVersion: ref.APIVersion,
	}
}

// PinnedObjectReference references a Kubernetes API object by UID, name and namespace.
// The reference is therefore linked to one specific object and will become invalid if the
// object is replaced by another one with the same name.
type PinnedObjectReference struct {
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
	// UID of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
	// +optional
	UID types.UID `json:"uid,omitempty" protobuf:"bytes,4,opt,name=uid,casttype=k8s.io/apimachinery/pkg/types.UID"`
}

func (r *PinnedObjectReference) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(r.APIVersion, r.Kind)
}

func (r *PinnedObjectReference) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	r.APIVersion, r.Kind = gvk.ToAPIVersionAndKind()
}

func PinnedObjectReferenceFromCore(ref corev1.ObjectReference) PinnedObjectReference {
	return PinnedObjectReference{
		Kind:       ref.Kind,
		Namespace:  ref.Namespace,
		Name:       ref.Name,
		APIVersion: ref.APIVersion,
		UID:        ref.UID,
	}
}

// Unpin removes the UID of the PinnedObjectReference turning it into
func (r PinnedObjectReference) Unpin() *ObjectReference {
	return &ObjectReference{
		Kind:       r.Kind,
		Name:       r.Name,
		APIVersion: r.APIVersion,
		Namespace:  r.Namespace,
	}
}

// LocalObjectReference references a Kubernetes API object only by name and can be used to reference objects in the same namespace.
type LocalObjectReference struct {
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
}

func (r *LocalObjectReference) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(r.APIVersion, r.Kind)
}

func (r *LocalObjectReference) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	r.APIVersion, r.Kind = gvk.ToAPIVersionAndKind()
}

// FullRef converts a LocalObjectReference to an ObjectReference by adding a namespace.
func (r LocalObjectReference) FullRef(namespace string) *ObjectReference {
	return &ObjectReference{
		Kind:       r.Kind,
		Name:       r.Name,
		APIVersion: r.APIVersion,
		Namespace:  namespace,
	}
}
