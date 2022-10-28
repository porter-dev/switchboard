package types

type Driver interface {
	PreApply(*Resource) error
	Apply(*Resource) error
	PostApply(*Resource) error
	OnError(*Resource, error)
}
