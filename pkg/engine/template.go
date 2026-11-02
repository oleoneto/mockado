package engine

var _ Engine = &Template{}

type Template struct{}

func (e *Template) Render(body []byte) ([]byte, error) { panic("not implemented") }
