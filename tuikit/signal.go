package tuikit

type SignalHandler func()

type Signal struct {
	handlers []SignalHandler
}

func (s *Signal) Attach(handler SignalHandler) int {
	for i, h := range s.handlers {
		if h == nil {
			s.handlers[i] = handler
			return i
		}
	}

	s.handlers = append(s.handlers, handler)
	return len(s.handlers) - 1
}

func (s *Signal) Detach(handle int) {
	s.handlers[handle] = nil
}

type SignalPublisher struct {
	signal Signal
}

func (p *SignalPublisher) Signal() *Signal {
	return &p.signal
}

func (p *SignalPublisher) Publish() {
	for _, handler := range p.signal.handlers {
		if handler != nil {
			handler()
		}
	}
}
