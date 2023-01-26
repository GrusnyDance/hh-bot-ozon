package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hh-bot-ozon/internal"
	"hh-bot-ozon/repository"
	"log"
	"os"
)

import (
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Println("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем время обновления
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // при открытии соединения оно живет 60 секунд

	// пул подключений к базе
	pool, err := repository.InitRep()
	if err != nil {
		log.Println(err)
	}
	defer pool.Close()

	instance := repository.Instance{Db: pool}
	checkDemon := internal.CheckUpdates(&instance)

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	for update := range updates {

		// универсальный ответ на любое сообщение
		reply := "I do not process any messages except commands"
		if update.Message == nil { // ignore any non-Message updates
			if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "this is meme dude")
				bot.Request(callback)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				bot.Send(msg)
				continue
			} else {
				continue
			}
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if !update.Message.IsCommand() {
			msg.Text = reply
		}

		siteUrl := os.Getenv("OZON_QUERY")

		if update.Message.IsCommand() {
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "list_of_commands":
				msg.Text = "I am command"
				k := []tgbotapi.InlineKeyboardButton{{Text: "lalala", CallbackData: &siteUrl}}
				Rm := tgbotapi.NewInlineKeyboardMarkup(k)
				msg.ReplyMarkup = Rm
			case "info":
				msg.Text = "I will show you recent internship positions"
			case "check_site":
				msg.Text = "Click the button"
				k := []tgbotapi.InlineKeyboardButton{{Text: "ozon", URL: &siteUrl}}
				Rm := tgbotapi.NewInlineKeyboardMarkup(k)
				msg.ReplyMarkup = Rm
			case "start":
				msg.Text = "This bot parses entry level vacancies from Ozon, tape /get_open_positions"
				go func() { // вынести в отдельного демона
					for {
						select {
						case newVac := <-checkDemon:
							notify := tgbotapi.NewMessage(update.Message.Chat.ID, "")
							notify.Text = newVac[0]
							k := []tgbotapi.InlineKeyboardButton{{Text: newVac[0], URL: &newVac[1]}}
							Rm := tgbotapi.NewInlineKeyboardMarkup(k)
							notify.ReplyMarkup = Rm
							bot.Send(notify)
						}
					}
				}()
			case "get_open_positions":
				vacMap := internal.GetOpenPositions()
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
			//case "number_of_users":
			//	msg.Text = "I calculate some stat of users"
			default:
				msg.Text = "I don't know that command"
			}
		}
		// отправляем
		bot.Send(msg)
	}
}
