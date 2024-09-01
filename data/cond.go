package data

import (
	"golang.org/x/exp/constraints"
)

// Qop - equivalent to ternary operator in other languages
func Qop[T any](cond bool, onTrue T, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

func NonEmpty[T comparable](vals ...T) T {
	var empty T
	for _, val := range vals {
		if val != empty {
			return val
		}
	}
	return empty
}

// CondExec - if cond is true run the trueFunc otherwise run the flase func
func CondExec(cond bool, trueFunc, falseFunc func() error) error {
	if cond {
		return trueFunc()
	}
	return falseFunc()
}

func OneOf[T comparable](tgt T, opts ...T) bool {
	for _, opt := range opts {
		if tgt == opt {
			return true
		}
	}
	return false
}

func Min[T constraints.Ordered](in ...T) T {
	var out T
	if len(in) == 0 {
		return out
	}

	out = in[0]
	for _, val := range in {
		if val < out {
			out = val
		}
	}

	return out
}

func Max[T constraints.Ordered](in ...T) T {
	var out T
	if len(in) == 0 {
		return out
	}

	out = in[0]
	for _, val := range in {
		if val > out {
			out = val
		}
	}

	return out
}
