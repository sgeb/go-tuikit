package binding

type BoolProperty interface {
	Get() bool
	Set(bool)
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

func (p *boolPropertyBase) Get() bool {
	return p.Property.Get().(bool)
}

func (p *boolPropertyBase) Set(v bool) {
	p.Property.Set(v)
}
