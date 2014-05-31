package databinding

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

func (sp *stringPropertyBase) Get() string {
	return sp.Property.Get().(string)
}

func (sp *stringPropertyBase) Set(s string) error {
	return sp.Property.Set(s)
}
