package binding

type ByteProperty interface {
	Get() byte
	Set(byte)
	Subscribe() <-chan struct{}
	Dispose()
}

type bytePropertyBase struct {
	Property
}

func NewByteProperty() ByteProperty {
	return &bytePropertyBase{
		NewProperty(),
	}
}
func (p *bytePropertyBase) Get() byte {
	return p.Property.Get().(byte)
}

func (p *bytePropertyBase) Set(v byte) {
	p.Property.Set(v)
}
