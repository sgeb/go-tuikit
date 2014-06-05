package binding

type BoolProperty interface {
	Get() bool
	Set(bool) error
	ReadOnly() bool
	Subscribe() <-chan struct{}
	Dispose()
}

type boolPropertyBase struct {
	Property
}

func NewBoolProperty() BoolProperty {
	return &boolPropertyBase{
		NewProperty(),
	}
}

func NewReadOnlyBoolProperty() BoolProperty {
	return &boolPropertyBase{
		NewReadOnlyProperty(),
	}
}

func (p *boolPropertyBase) Get() bool {
	return p.Property.Get().(bool)
}

func (p *boolPropertyBase) Set(s bool) error {
	return p.Property.Set(s)
}
