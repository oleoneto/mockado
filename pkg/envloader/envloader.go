package envloader

import (
	"os"
	"strings"

	"github.com/oleoneto/mockado/pkg/logger"
	"github.com/oleoneto/mockado/pkg/redactor"
	"github.com/sirupsen/logrus"
)

type EnvLoader struct {
	redactor redactor.Redactor
	logger   *logrus.Logger
}

type EnvLoaderConfig func(*EnvLoader)

type EnvVar struct{ Name, Value string }

func WithLogger(logger *logrus.Logger) EnvLoaderConfig {
	return func(el *EnvLoader) { el.logger = logger }
}

func WithRedactor(r *redactor.R) EnvLoaderConfig {
	return func(el *EnvLoader) { el.redactor = r }
}

func (l *EnvLoader) LoadVar(raw string) string {
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

	// If key is a sensitive var, don't load it. Simply return the corresponding `redacted pattern`.
	if s, matched := l.redactor.RedactMatched(key); matched {
		return s
	}

	v := os.Getenv(key)
	return l.redactor.Redact(v)
}

func (l *EnvLoader) Environment() map[string]string {
	variables := make(map[string]string)

	for _, pair := range os.Environ() {
		a := strings.SplitN(pair, "=", 2)
		key, value := a[0], a[1]

		if l.redactor != nil {
			// Check if key should be included at all
			if _, matched := l.redactor.RedactMatched(key); matched {
				continue
			}

			// Process value for possible redactions
			variables[key] = l.redactor.Redact(value)
			continue
		}

		variables[key] = value
	}

	return variables
}

func New(options ...EnvLoaderConfig) *EnvLoader {
	loader := &EnvLoader{
		redactor: noopRedactor{},
		logger:   logger.GetLogger(),
	}

	for _, fn := range options {
		fn(loader)
	}

	return loader
}

// noopRedactor is a redactor that never matches
type noopRedactor struct{}

func (noopRedactor) Redact(s string) string                { return s }
func (noopRedactor) RedactMatched(s string) (string, bool) { return s, false }
