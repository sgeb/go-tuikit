package binding

type Uint32Property interface {
	Get() uint32
	Set(uint32)
	Subscribe() <-chan struct{}
	Dispose()
}

type uint32PropertyBase struct {
	Property
}

func NewUint32Property() Uint32Property {
	return &uint32PropertyBase{
		NewProperty(),
	}
}

func (p *uint32PropertyBase) Get() uint32 {
	return p.Property.Get().(uint32)
}

func (p *uint32PropertyBase) Set(v uint32) {
	p.Property.Set(v)
}
