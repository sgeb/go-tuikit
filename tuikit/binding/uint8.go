package binding

type Uint8Property interface {
	Get() uint8
	Set(uint8)
	Subscribe() <-chan struct{}
	Dispose()
}

type uint8PropertyBase struct {
	Property
}

func NewUint8Property() Uint8Property {
	return &uint8PropertyBase{
		NewProperty(),
	}
}

func (p *uint8PropertyBase) Get() uint8 {
	return p.Property.Get().(uint8)
}

func (p *uint8PropertyBase) Set(v uint8) {
	p.Property.Set(v)
}
