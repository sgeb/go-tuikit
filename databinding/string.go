package databinding

type StringProperty interface {
	Get() string
	Set(string) error
	ReadOnly() bool
	Watch() WatcherChan
	Dispose()
}

type stringPropertyBase struct {
	readOnly bool
	val      string
}

func newString(readOnly bool, val string) StringProperty {
	return &stringPropertyBase{
		readOnly: readOnly,
		val:      val,
	}
}

func NewStringProperty(val string) StringProperty {
	return newString(false, val)
}

func NewReadOnlyStringProperty(val string) StringProperty {
	return newString(true, val)
}

func (sb *stringPropertyBase) Get() string {
	return sb.val
}

func (sb *stringPropertyBase) Set(s string) error {
	if sb.readOnly {
		return ErrPropertyReadOnly
	}

	sb.val = s
	return nil
}

func (sb *stringPropertyBase) ReadOnly() bool {
	return sb.readOnly
}

func (sb *stringPropertyBase) Watch() WatcherChan {
	return nil
}

func (sb *stringPropertyBase) Dispose() {
	return
}
