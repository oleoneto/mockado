package engine

type Engine interface {
	Render(body []byte) ([]byte, error)
}
