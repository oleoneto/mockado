package helpers

import (
	"bytes"
	"encoding/json"
)

// JSONPrettyPrint - produces a prettified/beautified JSON string (with proper spacing and indentation).
func JSONPrettyPrint(raw string) string {
	var output bytes.Buffer

	indentPrefix := ""
	indentTokens := "\t"
	err := json.Indent(&output, []byte(raw), indentPrefix, indentTokens)
	if err != nil {
		return raw
	}

	return output.String()
}
