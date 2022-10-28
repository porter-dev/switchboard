package types

type Driver interface {
	PreApply(*YAMLNode[*Resource]) error
	Apply(*YAMLNode[*Resource]) error
	PostApply(*YAMLNode[*Resource]) error
	OnError(*YAMLNode[*Resource], error)
}
