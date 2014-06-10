package binding

type Float32Property interface {
	Get() float32
	Set(float32)
	Subscribe() <-chan struct{}
	Dispose()
}

type float32PropertyBase struct {
	Property
}

func NewFloat32Property() Float32Property {
	return &float32PropertyBase{
		NewProperty(),
	}
}

func (p *float32PropertyBase) Get() float32 {
	return p.Property.Get().(float32)
}

func (p *float32PropertyBase) Set(v float32) {
	p.Property.Set(v)
}
