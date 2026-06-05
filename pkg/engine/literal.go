package engine

var _ Engine = &Literal{}

type Literal struct{}

func (e *Literal) Render(body []byte) ([]byte, error) { return body, nil }
