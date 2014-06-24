package binding

type RuneProperty interface {
	Get() rune
	Set(rune)
	Subscribe() <-chan struct{}
	Dispose()
}

type runePropertyBase struct {
	Property
}

func NewRuneProperty() RuneProperty {
	return &runePropertyBase{
		NewProperty(),
	}
}

func (p *runePropertyBase) Get() rune {
	return p.Property.Get().(rune)
}

func (p *runePropertyBase) Set(v rune) {
	p.Property.Set(v)
}
