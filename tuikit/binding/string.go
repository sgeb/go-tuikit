package binding

type StringProperty interface {
	Get() string
	Set(string) error
	ReadOnly() bool
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

func NewReadOnlyStringProperty() StringProperty {
	return &stringPropertyBase{
		NewReadOnlyProperty(),
	}
}

func (p *stringPropertyBase) Get() string {
	return p.Property.Get().(string)
}

func (p *stringPropertyBase) Set(s string) error {
	return p.Property.Set(s)
}
