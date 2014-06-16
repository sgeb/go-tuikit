package binding

type Uint16Property interface {
	Get() uint16
	Set(uint16)
	Subscribe() <-chan struct{}
	Dispose()
}

type uint16PropertyBase struct {
	Property
}

func NewUint16Property() Uint16Property {
	return &uint16PropertyBase{
		NewProperty(),
	}
}

func (p *uint16PropertyBase) Get() uint16 {
	return p.Property.Get().(uint16)
}

func (p *uint16PropertyBase) Set(v uint16) {
	p.Property.Set(v)
}
