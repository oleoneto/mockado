package redactor

import "regexp"

const RedactedPattern string = "****"

type SecretPattern struct {
	Name    string `yaml:"name" json:"name"`
	Pattern []byte `yaml:"pattern" json:"pattern"`
}

// DefaultPatterns is the collection of high-confidence sensitive string detectors
func DefaultPatterns() []SecretPattern {
	return []SecretPattern{
		// --- Version Control & Git Platforms ---
		{
			Name:    "GitHub Fine-grained PAT",
			Pattern: []byte(`github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}`),
		},
		{
			Name:    "GitHub Classic PAT",
			Pattern: []byte(`ghp_[a-zA-Z0-9]{36}`),
		},
		{
			Name:    "GitHub OAuth Access Token",
			Pattern: []byte(`gho_[a-zA-Z0-9]{36}`),
		},
		{
			Name:    "GitLab Personal Access Token",
			Pattern: []byte(`glpat-[a-zA-Z0-9\-]{20}`),
		},

		// --- Cloud Providers ---
		{
			Name:    "AWS Access Key ID",
			Pattern: []byte(`\b(AKIA|ASIA|AGPA|AIDA|AROA|AIPA)[A-Z0-9]{16}\b`),
		},
		{
			Name:    "AWS Secret Access Key",
			Pattern: []byte(`(?i)aws(.{0,20})?['\"][0-9a-zA-Z\/+]{40}['\"]`),
		},
		{
			Name:    "GCP API Key",
			Pattern: []byte(`\bAIza[0-9A-Za-z-_]{35}\b`),
		},

		// --- AI & Developer Services ---
		{
			Name:    "OpenAI API Key",
			Pattern: []byte(`sk-(?:proj-|svcacct-|live-|test-)?[a-zA-Z0-9\-_]{80,}`),
		},
		{
			Name:    "Anthropic API Key",
			Pattern: []byte(`sk-ant-sid01-[a-zA-Z0-9_\-]{93}`),
		},
		{
			Name:    "Stripe Live Secret Key",
			Pattern: []byte(`sk_live_[0-9a-zA-Z]{24,}`),
		},
		{
			Name:    "Slack Token",
			Pattern: []byte(`xox[bapr]-[0-9]{12}-[0-9]{12}-[a-zA-Z0-9]{24}`),
		},

		// --- Generic & Cryptographic ---
		{
			Name:    "Asymmetric Private Key",
			Pattern: []byte(`-----BEGIN (RSA|DSA|EC|OPENSSH|PGP)? PRIVATE KEY-----`),
		},
		{
			Name:    "JSON Web Token (JWT)",
			Pattern: []byte(`\beyj[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.?[A-Za-z0-9_-]*\b`),
		},
		{
			Name:    "Generic Password/Key Assignment",
			Pattern: []byte(`(?i)\b(api[_\s-]?key|secret|password|passwd|token)\b\s*[:=]\s*['\"][A-Za-z0-9_\-]{16,}['\"]`),
		},
	}
}

func compilePatterns(patterns ...SecretPattern) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(string(p.Pattern))
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}

func (s *SecretPattern) UnmarshalYAML(unmarshal func(any) error) error {
	type s2 SecretPattern

	s3 := &struct {
		Pattern string `yaml:"pattern"`
		*s2
	}{
		s2: (*s2)(s),
	}

	if err := unmarshal(&s3); err != nil {
		return err
	}

	s.Pattern = []byte(s3.Pattern)
	return nil
}
