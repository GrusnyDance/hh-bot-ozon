package internal

import (
	"fmt"
	"hh-bot-ozon/repository"
	"time"
)

func CheckUpdates(i *repository.Instance, ch chan []string) {
	for { //собираем данные
		fmt.Println("I work")
		mapVac := GetOpenPositions()
		if len(*mapVac) > 0 {
			fmt.Println("i am line 16")
			if resMap, ok := i.GetVac(mapVac); ok {
				for k, v := range *resMap {
					ret := []string{k, v}
					ch <- ret
				}
			}
		}
		time.Sleep(3 * time.Minute)
	}
}
