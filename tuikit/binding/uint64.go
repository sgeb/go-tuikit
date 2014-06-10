package binding

type Uint64Property interface {
	Get() uint64
	Set(uint64)
	Subscribe() <-chan struct{}
	Dispose()
}

type uint64PropertyBase struct {
	Property
}

func NewUint64Property() Uint64Property {
	return &uint64PropertyBase{
		NewProperty(),
	}
}

func (p *uint64PropertyBase) Get() uint64 {
	return p.Property.Get().(uint64)
}

func (p *uint64PropertyBase) Set(v uint64) {
	p.Property.Set(v)
}
