package internal

import (
	"time"
)

func Check_updates() <-chan string { //функция вернет канал типа string
	sendVac := make(chan string)
	go func() {
		for { //собираем данные
			mapVac := GetOpenPositions()
			var newVacs bool
			if len(*mapVac) > 0 {
				// code to compare map with db
			}
			if newVacs {
				sendVac <- "LALALA" // some ret from line 33
			}
			time.Sleep(3 * time.Minute)
		}
	}()
	return sendVac
}
