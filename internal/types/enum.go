package types

type Kind uint8

const (
	ConfigMap Kind = iota
	Secret
	Namespace
)
