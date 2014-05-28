package databinding

import "errors"

var (
	ErrPropertyReadOnly error = errors.New("read-only property")
)

type WatcherChan <-chan struct{}

type Property interface {
	Get() interface{}
	Set(interface{}) error
	ReadOnly() bool
	Watch() WatcherChan
	Dispose()
}

type propertyBase struct {
	readOnly bool
	val      interface{}
}

func newProperty(readOnly bool) Property {
	return &propertyBase{
		readOnly: readOnly,
	}
}

func NewProperty() Property {
	return newProperty(false)
}

func NewReadOnlyProperty() Property {
	return newProperty(true)
}

func (p *propertyBase) Get() interface{} {
	return p.val
}

func (p *propertyBase) Set(nv interface{}) error {
	if p.ReadOnly() {
		return ErrPropertyReadOnly
	}

	p.val = nv
	return nil
}

func (p *propertyBase) ReadOnly() bool {
	return p.readOnly
}

func (p *propertyBase) Watch() WatcherChan {
	return nil
}

func (p *propertyBase) Dispose() {
	return
}
