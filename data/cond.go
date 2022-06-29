package data

// Qop - equivalent to ternary operator in other languages
func Qop[T any](cond bool, onTrue T, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

// CondExec - if cond is true run the trueFunc otherwise run the flase func
func CondExec(cond bool, trueFunc, falseFunc func() error) error {
	if cond {
		return trueFunc()
	}
	return falseFunc()
}
