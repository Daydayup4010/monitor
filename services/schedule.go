package services

import (
	"time"
)

func StartBuffFullUpdateScheduler() {
	go UpdateBuffFullData()
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateBuffFullData()
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

func StartUUFullUpdateScheduler() {
	go UpdateUUFullData()
	ticker := time.NewTicker(8 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateUUFullData()
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
