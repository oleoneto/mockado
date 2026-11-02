package redactor

import (
	"bytes"
	"regexp"

	"github.com/oleoneto/mockado/pkg/helpers"
	"github.com/oleoneto/mockado/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Redactor interface {
	Redact(value string) string
	RedactMatched(s string) (substitute string, matched bool)
}

type R struct {
	patterns        []*regexp.Regexp
	logger          *logrus.Logger
	redactedPattern *string
}

type RedactorConfig func(*R)

var redactor *R = &R{}

// Redact redacts the input string if it matches any of the redactor's configured patterns.
//
// It returns a substitute string or the original string if no patterns were matched.
func (r *R) Redact(value string) string {
	substitute, _ := r.redactValue(value)
	return substitute
}

// RedactMatched redacts the input string if it matches any of the redactor's configured patterns.
//
// If returns a substitute string (which could be the original string) and a boolean parameter indicating if a match was found.
func (r *R) RedactMatched(s string) (substitute string, matched bool) {
	return r.redactValue(s)
}

func (r *R) redactValue(s string) (substitution string, matched bool) {
	for _, blocked := range r.patterns {
		if blocked.MatchString(s) {
			bf := bytes.Buffer{}

			bf.Write([]byte(s))

			return *r.redactedPattern, true
		}
	}

	return s, false
}

func (r *R) addToList(blocked ...[]byte) *R {
	if len(blocked) == 0 {
		return r
	}

	patterns := helpers.Map(blocked, func(_ int, x []byte) SecretPattern {
		return SecretPattern{Pattern: x}
	})

	compiled, err := compilePatterns(patterns...)
	if err != nil {
		panic(err)
	}

	r.patterns = append(r.patterns, compiled...)

	return r
}

func WithLogger(logger *logrus.Logger) RedactorConfig {
	return func(r *R) { r.logger = logger }
}

func WithRedactedPattern(s string) RedactorConfig {
	return func(r *R) { r.redactedPattern = new(s) }
}

func WithBlockedList(blocked ...[]byte) RedactorConfig {
	return func(r *R) { redactor.addToList(blocked...) }
}

func WithBlockedPatterns(patterns ...SecretPattern) RedactorConfig {
	return func(r *R) {
		if len(patterns) == 0 {
			return
		}

		blocked := helpers.Map(patterns, func(_ int, sp SecretPattern) []byte { return sp.Pattern })
		redactor.addToList(blocked...)
	}
}

func WithDefaultPatterns() RedactorConfig {
	return func(r *R) {
		patterns, err := compilePatterns(DefaultPatterns()...)
		if err != nil {
			panic(err)
		}

		r.patterns = append(r.patterns, patterns...)
	}
}

func New(options ...RedactorConfig) *R {
	redactor = &R{
		logger:          logger.GetLogger(),
		redactedPattern: new(RedactedPattern),
	}

	for _, fn := range options {
		fn(redactor)
	}

	redactor.logger.
		WithField("package", "redactor").
		Tracef("Added %d patterns\n", len(redactor.patterns))

	return redactor
}
