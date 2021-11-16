Drivers satisfy a common interface. 

```go
type Driver interface {
  ShouldApply(resource *Resource) bool
  Apply(resource *Resource) (*Resource, error)
  Destroy(resource *Resource) (*Resource, error)
}
```