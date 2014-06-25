package binding

type UintProperty interface {
	Get() uint
	Set(uint)
	Subscribe() <-chan struct{}
	Dispose()
}

type uintPropertyBase struct {
	Property
}

func NewUintProperty() UintProperty {
	return &uintPropertyBase{
		NewProperty(),
	}
}

func (p *uintPropertyBase) Get() uint {
	return p.Property.Get().(uint)
}

func (p *uintPropertyBase) Set(v uint) {
	p.Property.Set(v)
}
