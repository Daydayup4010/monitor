package services

import (
	"uu/config"
)

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				config.Log.Panicf("goroutine panic recovered: %v", r)
			}
		}()
		fn()
	}()
}
