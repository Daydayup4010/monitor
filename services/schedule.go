package services

import (
	"time"
)

func StartFullUpdateScheduler() {
	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateFullData()
	}
}

func StartVerifyToken() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go func() {
			VerifyBuffToken()
			VerifyUUToken()
		}()
	}
}
