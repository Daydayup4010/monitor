package services

import (
	"time"
)

func StartFullUpdateScheduler() {
	ticker := time.NewTicker(25 * time.Minute)
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

//
//func StartUpdateSkin() {
//	ticker := time.NewTicker(1 * time.Minute)
//	defer ticker.Stop()
//	for range ticker.C {
//		go models.UpdateSkinItems()
//	}
//}
