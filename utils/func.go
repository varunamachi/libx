package utils

import "github.com/varunamachi/libx/errx"

type Pair[T1, T2 any] struct {
	One T1
	Two T2
}

func N2[T1, T2 any](a1 T1, a2 T2) *Pair[T1, T2] {
	return &Pair[T1, T2]{
		One: a1,
		Two: a2,
	}
}

type Triad[T1, T2, T3 any] struct {
	One   T1
	Two   T2
	Three T3
}

func N3[T1, T2, T3 any](a1 T1, a2 T2, a3 T3) *Triad[T1, T2, T3] {
	return &Triad[T1, T2, T3]{
		One:   a1,
		Two:   a2,
		Three: a3,
	}
}

func Call(funcs ...func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func CallWith1Arg[T any](fn func(arg T) error, args ...T) error {
	for _, arg := range args {
		if err := fn(arg); err != nil {
			return errx.Errf(err, "failed to call fn of %v", arg)
		}
	}
	return nil
}

func CallWith2Args[T1, T2 any](
	fn func(arg1 T1, arg2 T2) error,
	args ...*Pair[T1, T2]) error {

	for _, arg := range args {
		if err := fn(arg.One, arg.Two); err != nil {
			return errx.Errf(err, "failed to call fn of %v", arg)
		}
	}
	return nil
}

func CallWith3Arg[T1, T2, T3 any](
	fn func(arg1 T1, arg2 T2, arg3 T3) error,
	args ...*Triad[T1, T2, T3]) error {

	for _, arg := range args {
		if err := fn(arg.One, arg.Two, arg.Three); err != nil {
			return errx.Errf(err, "failed to call fn of %v", arg)
		}
	}
	return nil
}
