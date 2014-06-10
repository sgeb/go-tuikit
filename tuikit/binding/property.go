package binding

type Property interface {
	Get() interface{}
	Set(interface{})
	Subscribe() <-chan struct{}
	Dispose()
}

type propertyBase struct {
	val           interface{}
	subscriptions []chan struct{}
}

func NewProperty() Property {
	return &propertyBase{}
}

func (p *propertyBase) Get() interface{} {
	return p.val
}

func (p *propertyBase) Set(v interface{}) {
	p.val = v

	for _, c := range p.subscriptions {
		c <- struct{}{}
	}
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
