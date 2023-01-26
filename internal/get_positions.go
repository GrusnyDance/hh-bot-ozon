package internal

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetOpenPositions() *map[string]string {
	url := os.Getenv("OZON_QUERY")
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var numOfVacancies int
	doc.Find("div.search").Each(func(i int, s *goquery.Selection) {
		str := s.Find("div.search__count").Text()
		if str != "" {
			numStr := strings.Fields(str)[1]
			numOfVacancies, _ = strconv.Atoi(numStr)
		}
	})

	vacMap := make(map[string]string, numOfVacancies)
	if numOfVacancies > 0 {
		doc.Find("div.finder__main").Find("div.results__items").Find("div.wr").Each(func(i int, s *goquery.Selection) {
			str := s.Find("h6.result__title").Text()
			strUrl, _ := s.Find("a").Attr("href")
			if str != "" {
				vacMap[strings.Trim(str, "\n ")] = os.Getenv("OZON_PREFIX") + strUrl
			}
		})
	}
	return &vacMap
}
