package databinding

type StringProperty interface {
	Get() string
	Set(string) error
	ReadOnly() bool
	Watch() WatcherChan
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

func (sp *stringPropertyBase) ReadOnly() bool {
	return sp.Property.ReadOnly()
}

func (sp *stringPropertyBase) Watch() WatcherChan {
	return sp.Property.Watch()
}

func (sp *stringPropertyBase) Dispose() {
	sp.Property.Dispose()
}
