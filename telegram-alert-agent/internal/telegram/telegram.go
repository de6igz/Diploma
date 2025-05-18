package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewTelegramBot инициализирует клиента Telegram Bot API.
func NewTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

// SendMessage отправляет сообщение в Telegram через переданного бота.
func SendMessage(bot *tgbotapi.BotAPI, chatID string, text string) error {
	// Преобразуем chatID в int64; ожидается, что в сообщении приходит числовой id (например, "456")
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(id, text)
	_, err = bot.Send(msg)
	return err
}
