package envloader_test

import (
	"os"
	"testing"

	"github.com/oleoneto/mockado/pkg/envloader"
	"github.com/oleoneto/mockado/pkg/redactor"
)

func Test_parseTemplateVariable(t *testing.T) {
	tests := []struct {
		name     string
		options  []redactor.RedactorConfig
		variable string
		want     string
	}{
		{
			name:     "0",
			variable: "{{ env.REPOSITORY }}",
			want:     "/home/root/myrepo",
			options:  []redactor.RedactorConfig{},
		},
		{
			name:     "env.ABC",
			variable: "{{ env.ABC }}",
			want:     "/var/abc",
			options:  []redactor.RedactorConfig{},
		},
		{
			name:     "env.XYZ",
			variable: "{{ env.XYZ }}",
			want:     "/tmp/random",
			options:  []redactor.RedactorConfig{},
		},
		{
			name:     "env.DB_PASSWORD",
			variable: "{{ env.DB_PASSWORD }}",
			options:  []redactor.RedactorConfig{redactor.WithBlockedList([]byte("DB|PASSWORD"))},
			want:     "****",
		},
		{
			name:     "env.USER_PASSWORD",
			variable: "{{ env.USER_PASSWORD }}",
			options:  []redactor.RedactorConfig{redactor.WithBlockedList([]byte("DB|PASSWORD"))},
			want:     "****",
		},
		{
			name:     "env.GIT_CREDS",
			variable: "{{ env.GIT_CREDS }}",
			options:  []redactor.RedactorConfig{redactor.WithBlockedList([]byte("CRED"))},
			want:     "****",
		},
		{
			name:     "env.ACCESS_KEY_ID",
			variable: "{{ env.ACCESS_KEY_ID }}", // Part of the default block list
			options: []redactor.RedactorConfig{
				redactor.WithBlockedList([]byte("KEY")),
				redactor.WithRedactedPattern("[sensitive]"),
			},
			want: "[sensitive]",
		},
	}

	// redactor.WithBlockedList([]byte("CRED"), []byte("DB")),

	for _, tt := range tests {
		os.Setenv("ABC", "/var/abc")
		os.Setenv("XYZ", "/tmp/random")
		os.Setenv("REPOSITORY", "/home/root/myrepo")

		r := redactor.New(tt.options...)

		loader := envloader.New(
			envloader.WithRedactor(r),
		)

		t.Run(tt.name, func(t *testing.T) {
			got := loader.LoadVar(tt.variable)
			if got != tt.want {
				t.Errorf("parseTemplateVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
