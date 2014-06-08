package binding

type ByteProperty interface {
	Get() byte
	Set(byte) error
	ReadOnly() bool
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

func NewReadOnlyByteProperty() ByteProperty {
	return &bytePropertyBase{
		NewReadOnlyProperty(),
	}
}

func (p *bytePropertyBase) Get() byte {
	return p.Property.Get().(byte)
}

func (p *bytePropertyBase) Set(v byte) error {
	return p.Property.Set(v)
}
