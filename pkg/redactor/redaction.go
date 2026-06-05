package redactor

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"regexp"
	"sync"

	"github.com/cleopatrio/mockado/pkg/helpers"
	"github.com/cleopatrio/mockado/pkg/logger"
)

var once sync.Once

var redactor *Redactor = &Redactor{}

type Redactor struct{ patterns []*regexp.Regexp }

func (r *Redactor) Redact(value string) string { out, _ := r.redact(value); return out }

func (r *Redactor) RedactKey(key string) (v string, matched bool) { return r.redact(key) }

func (r *Redactor) redact(s string) (out string, matched bool) {
	for _, blocked := range r.patterns {
		if blocked.MatchString(s) {
			bf, hash := bytes.Buffer{}, sha1.New()

			bf.Write([]byte(s))
			io.Copy(hash, &bf)

			log.
				WithField("package", "redaction").
				Debugf("Value %s redacted\n", hex.EncodeToString(hash.Sum(nil)))

			return RedactedPattern, true
		}
	}

	return s, false
}

func (r *Redactor) addToList(blocked ...[]byte) *Redactor {
	patterns := helpers.Map(blocked, func(_ int, x []byte) SecretPattern {
		return SecretPattern{Pattern: x}
	})

	compiled, err := compilePatterns(patterns...)
	if err != nil {
		log.Fatalln(err)
	}

	r.patterns = append(r.patterns, compiled...)

	return r
}

func WithBlockList(blocked ...[]byte) *Redactor { return redactor.addToList(blocked...) }

func New() *Redactor {
	once.Do(func() {
		patterns, err := compilePatterns(DefaultPatterns()...)
		if err != nil {
			log.Panic(err)
		}

		log.
			WithField("package", "redactor").
			Tracef("Called New() with %d patterns\n", len(patterns))

		redactor = &Redactor{patterns: patterns}
	})

	return redactor
}

var log = logger.GetLogger("redactor")
