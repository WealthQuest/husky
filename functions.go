package husky

import "runtime/debug"

func ToPoint[P any](p P) *P {
	return &p
}

func If[R any](logic bool, r1, r2 R) R {
	if logic {
		return r1
	} else {
		return r2
	}
}

func Iff[R any](logic bool, f1, f2 func() R) R {
	if logic {
		return f1()
	} else {
		return f2()
	}
}

func NilDefault[P any](data *P, defaultVal P) P {
	if data == nil {
		return defaultVal
	} else {
		return *data
	}
}

func Filter[S ~[]E, E any](datas S, f func(item E) bool) S {
	result := make(S, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		if r {
			result = append(result, data)
		}
	}
	return result
}

func IsMatch[T comparable](v T, matchers ...T) bool {
	for i := range matchers {
		if v == matchers[i] {
			return true
		}
	}
	return false
}

func Map[P any, R any](datas []P, f func(item P) R) []R {
	result := make([]R, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		result = append(result, r)
	}
	return result
}

func Convert[D any, R any](d D, call func(d D) R) R {
	return call(d)
}

func Go(f func()) {
	go func() {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			println(err)
			debug.PrintStack()
		}()
		f()
	}()
}
