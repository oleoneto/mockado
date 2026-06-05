package engine

import (
	"fmt"
	"log"

	"github.com/google/go-jsonnet"
)

var _ Engine = &Jsonnet{}

type Jsonnet struct{}

func (e *Jsonnet) Render(body []byte) ([]byte, error) {
	vm := jsonnet.MakeVM()

	snippet := string(body)

	str, err := vm.EvaluateAnonymousSnippet("mockado.jsonnet", snippet)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(str)

	return []byte(str), nil
}
