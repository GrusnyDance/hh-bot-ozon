package internal

import (
	"hh-bot-ozon/repository"
	"time"
)

func CheckUpdates(i *repository.Instance) <-chan []string { //функция вернет канал типа string
	sendVac := make(chan []string)
	go func() {
		for { //собираем данные
			mapVac := GetOpenPositions()
			if len(*mapVac) > 0 {
				if resMap, ok := i.GetVac(mapVac); ok {
					for k, v := range *resMap {
						ret := []string{k, v}
						sendVac <- ret
					}
				}
			}
			time.Sleep(3 * time.Minute)
		}
	}()
	return sendVac
}
