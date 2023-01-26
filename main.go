package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

import (
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем время обновления
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0) // зачем
	u.Timeout = 60             // при открытии соединения оно живет 60 секунд

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := "I do not process any messages except commands"
		if update.Message == nil { // ignore any non-Message updates
			if update.CallbackQuery != nil {
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				// мы получили CallBack, обработали информацию, которую он нам нес,
				// и сообщили телеграмму, что мы обработали этот CallBack
				// в середине экрана черное окошко с текстом
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "this is meme dude")
				bot.Request(callback)

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				bot.Send(msg)
				continue
			} else {
				continue
			}
		}

		// логируем от кого какое сообщение пришло
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if !update.Message.IsCommand() {
			msg.Text = reply
		}
		test_url := "https://core.telegram.org/bots/api#replykeyboardmarkup"
		test_string := "hahaha"

		if update.Message.IsCommand() {
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "list_of_commands":
				msg.Text = "I am command"
				k := []tgbotapi.InlineKeyboardButton{{Text: "lalala", CallbackData: &test_string}}
				Rm := tgbotapi.NewInlineKeyboardMarkup(k)
				msg.ReplyMarkup = Rm
			case "get_open_positions":
				msg.Text = "I will show you recent internship positions"
			case "link":
				msg.Text = "I am link"
				k := []tgbotapi.InlineKeyboardButton{{Text: "tap me", URL: &test_url}, {Text: "haha", URL: &test_url}}
				Rm := tgbotapi.NewInlineKeyboardMarkup(k)
				msg.ReplyMarkup = Rm
			case "start":
				msg.Text = "This bot parses entry level vacancies from Ozon on hh"
			case "test":
				vacMap := getOpenPositions()
				var rows []tgbotapi.InlineKeyboardButton
				if len(*vacMap) > 0 {
					for key, value := range *vacMap {
						fmt.Println("key map is", key)
						msg.Text += key + "\n"
						rows = append(rows, tgbotapi.InlineKeyboardButton{Text: key, URL: &value})
					}
					Rm := tgbotapi.NewInlineKeyboardMarkup(rows)
					msg.ReplyMarkup = Rm
				} else {
					msg.Text = "No vacancies available"
				}
			case "number_of_users":
				msg.Text = "I calculate some stat of users" // available for me only
				// запрос в базу с пользователями
			default:
				msg.Text = "I don't know that command"
			}
		}

		// отправляем
		bot.Send(msg)
	}
}

func getOpenPositions() *map[string]string {
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
