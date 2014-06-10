package binding

type StringProperty interface {
	Get() string
	Set(string)
	Subscribe() <-chan struct{}
	Dispose()
}

type stringPropertyBase struct {
	Property
}

func NewStringProperty() StringProperty {
	return &stringPropertyBase{
		NewProperty(),
	}
}

func (p *stringPropertyBase) Get() string {
	return p.Property.Get().(string)
}

func (p *stringPropertyBase) Set(v string) {
	p.Property.Set(v)
}
