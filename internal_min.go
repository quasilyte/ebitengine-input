package input

import "golang.org/x/exp/constraints"

// minOf is a placeholder implementation until Go 1.21's builtin min implementation is available.
func minOf[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
