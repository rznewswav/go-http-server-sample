package utils

func Repeat(count int, callback *func(*func())) {
	for i := 0; i < count; i++ {
		shouldCancel := false
		cancelHandler := func() {
			shouldCancel = true
		}
		(*callback)(&cancelHandler)
		if shouldCancel {
			return
		}
	}
}
