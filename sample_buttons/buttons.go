package sample_buttons

//
//internal main
//
//import (
//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//"log"
//"os"
//)
//
//import (
//"github.com/joho/godotenv"
//)
//
//var mainMenu = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("🏠 Главная"),
//		tgbotapi.NewKeyboardButton("🗒 Запись"),
//	),
//)
//
//func main() {
//	err := godotenv.Load() // 👈 load .env file
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	// Устанавливаем время обновления
//	// u - структура с конфигом для получения апдейтов
//	u := tgbotapi.NewUpdate(0) // зачем
//	u.Timeout = 60             // при открытии соединения оно живет 60 секунд
//
//	// используя конфиг u создаем канал в который будут прилетать новые сообщения
//	updates := bot.GetUpdatesChan(u)
//
//	// в канал updates прилетают структуры типа Update
//	// вычитываем их и обрабатываем
//	for update := range updates {
//		// универсальный ответ на любое сообщение
//		reply := "I do not process any messages except commands"
//		if update.Message == nil { // ignore any non-Message updates
//			if update.CallbackQuery != nil {
//				// Respond to the callback query, telling Telegram to show the user
//				// a message with the data received.
//				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
//				if _, err := bot.Request(callback); err != nil {
//					panic(err)
//				}
//
//				// And finally, send a message containing the data received.
//				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
//				if _, err := bot.Send(msg); err != nil {
//					panic(err)
//				}
//			} else {
//				continue
//			}
//		}
//
//		// логируем от кого какое сообщение пришло
//		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//		// создаем ответное сообщение
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
//		if !update.Message.IsCommand() {
//			msg.Text = reply
//		}
//
//		if update.Message.IsCommand() {
//			// Extract the command from the Message.
//			switch update.Message.Command() {
//			case "list_of_commands":
//				msg.Text = "I am command"
//				msg.ReplyMarkup = mainMenu
//			case "check_recent_vacancies":
//				msg.Text = "I will show you recent internship positions"
//			case "start":
//				msg.Text = "This bot parses entry level vacancies from Ozon on hh"
//			case "info":
//				msg.Text = "Some instructions for user"
//			case "number_of_users":
//				msg.Text = "I calculate some stat of users" // available for me only
//				// запрос в базу с пользователями
//			default:
//				msg.Text = "I don't know that command"
//			}
//		}
//
//		// отправляем
//		bot.Send(msg)
//	}
//}
