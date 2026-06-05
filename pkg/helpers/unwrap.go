package helpers

// UnwrapOrDefault - returns the value of the optional argument or a fallback value if the optional is nil.
func UnwrapOrDefault[T any](optional *T, fallback T) T {
	if optional != nil {
		return *optional
	}

	return fallback
}
