package envloader

import (
	"os"
	"strings"
	"sync"

	"github.com/cleopatrio/mockado/pkg/logger"
	"github.com/cleopatrio/mockado/pkg/redactor"
)

var once sync.Once

var loader *Loader

var log = logger.GetLogger("envloader")

type Loader struct{ redactor *redactor.Redactor }

type EnvVar struct{ Name, Value string }

func WithRedactor(redactor *redactor.Redactor) *Loader {
	loader = New()
	loader.redactor = redactor
	return loader
}

func (l *Loader) LoadVar(raw string) string {
	if raw == "" {
		return raw
	}

	raw = strings.TrimFunc(raw, func(r rune) bool {
		switch r {
		case ' ', '{', '}':
			return true
		}
		return false
	})

	x := strings.SplitN(raw, ".", 2)
	if len(x) != 2 {
		return raw
	}

	if x[0] != "env" {
		return raw
	}

	key := x[1]

	if k, redacted := loader.redactor.RedactKey(key); redacted {
		return k
	}

	y := os.Getenv(key)

	return loader.redactor.Redact(y)
}

func New() *Loader {
	once.Do(func() {
		loader = &Loader{}
		log.
			WithField("package", "envloader").
			Tracef("Called New()")
	})

	return loader
}
