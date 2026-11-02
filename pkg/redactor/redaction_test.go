package redactor_test

import (
	"testing"

	"github.com/oleoneto/mockado/pkg/redactor"
)

func TestRedactor_RedactMatched(t *testing.T) {
	tests := []struct {
		name           string
		options        []redactor.RedactorConfig
		input          string
		wantSubstitute string
		wantMatched    bool
	}{
		{
			name: "Live Tokens",
			options: []redactor.RedactorConfig{
				redactor.WithBlockedPatterns(
					redactor.SecretPattern{Name: "Stripe Key", Pattern: []byte(`sk_live_[0-9a-zA-Z]{12,}`)},
				),
			},
			input:          "sk_test_1sK3S4h347p2",
			wantSubstitute: "sk_test_1sK3S4h347p2",
			wantMatched:    false,
		},
		{
			name: "Live Tokens",
			options: []redactor.RedactorConfig{
				redactor.WithRedactedPattern("[redacted]"),
				redactor.WithBlockedPatterns(
					redactor.SecretPattern{Name: "Stripe Key", Pattern: []byte(`sk_fake_[0-9a-zA-Z]{12,}`)},
				),
			},
			input:          "sk_fake_000000000000",
			wantSubstitute: "[redacted]",
			wantMatched:    true,
		},
		{
			name: "Live Tokens",
			options: []redactor.RedactorConfig{
				redactor.WithRedactedPattern("****"),
				redactor.WithBlockedPatterns(
					redactor.SecretPattern{Name: "Stripe Key", Pattern: []byte(`sk_fake_[0-9a-zA-Z]{12,}`)},
				),
			},
			input:          "sk_fake_000000000000",
			wantSubstitute: "****",
			wantMatched:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := redactor.New(tt.options...)
			substitute, matched := r.RedactMatched(tt.input)
			if matched != tt.wantMatched {
				t.Errorf("RedactMatched() = %v, want %v", matched, tt.wantMatched)
			}
			if substitute != tt.wantSubstitute {
				t.Errorf("RedactMatched() = %v, want %v", substitute, tt.wantSubstitute)
			}
		})
	}
}
