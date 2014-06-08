package binding

import "errors"

var (
	errPropertyReadOnly error = errors.New("read-only property")
)

type Property interface {
	Get() interface{}
	Set(interface{}) error
	ReadOnly() bool
	Subscribe() <-chan struct{}
	Dispose()
}

type propertyBase struct {
	readOnly      bool
	val           interface{}
	subscriptions []chan struct{}
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

func (p *propertyBase) Set(v interface{}) error {
	if p.ReadOnly() {
		return errPropertyReadOnly
	}

	p.val = v

	for _, c := range p.subscriptions {
		c <- struct{}{}
	}

	return nil
}

func (p *propertyBase) ReadOnly() bool {
	return p.readOnly
}

func (p *propertyBase) Subscribe() <-chan struct{} {
	c := make(chan struct{}, 1)
	p.subscriptions = append(p.subscriptions, c)
	return c
}

func (p *propertyBase) Dispose() {
	for _, c := range p.subscriptions {
		close(c)
	}
}
