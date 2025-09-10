package services

import (
	"time"
)

func StartFullUpdateScheduler() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	go UpdateFullData()
	for range ticker.C {
		go UpdateFullData()
	}
}

func StartVerifyToken() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	go func() {
		VerifyBuffToken()
		VerifyUUToken()
	}()
	for range ticker.C {
		go func() {
			VerifyBuffToken()
			VerifyUUToken()
		}()
	}
}
