package engine

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var _ Engine = &Go{}

type Go struct{}

func (e *Go) Render(body []byte) ([]byte, error) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	_, err := i.Eval(string(body))
	if err != nil {
		return nil, err
	}

	v, err := i.Eval("router.Route")
	if err != nil {
		return nil, err
	}

	bodyFunc := v.Interface().(func() string)
	bodyOut := bodyFunc()

	return []byte(bodyOut), nil
}
