package types

// K8sResourceKind is the resource kind of kubernetes
type K8sResourceKind uint8

const (
	// ConfigMap is the kubernetes Resource that kind is configmap
	ConfigMap K8sResourceKind = iota
	// Secret is the kubernetes Resource that kind is secret
	Secret
	// Namespace is the namespace of kubernetes
	Namespace
)
