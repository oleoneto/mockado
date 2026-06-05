package helpers

// PointerTo - returns a reference to value.
func PointerTo[T any](value T) *T { return &value }
