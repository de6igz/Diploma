package telegram_repository

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramRepository описывает интерфейс для работы с Telegram.
type TelegramRepository interface {
	SendMessage(chatID string, text string) error
	// StartCommandListener запускает прослушивание входящих обновлений (команд) бота.
	StartCommandListener()
}

type telegramRepository struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramRepository создаёт экземпляр TelegramRepository, инициализируя бота.
func NewTelegramRepository(token string) (TelegramRepository, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &telegramRepository{bot: bot}, nil
}

// SendMessage отправляет сообщение в Telegram.
func (r *telegramRepository) SendMessage(chatID string, text string) error {
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.Entities = []tgbotapi.MessageEntity{{
		Type:   "pre", // Используем pre, чтобы Telegram показал зелёный блок
		Offset: 0,
		Length: len(text),
	}}
	_, err = r.bot.Send(msg)
	return err
}

// StartCommandListener запускает цикл обработки входящих обновлений от Telegram.
// Если пользователь отправляет команду /get_chat_id, бот отвечает сообщением с его chat id.
func (r *telegramRepository) StartCommandListener() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "get_chat_id":
				chatID := update.Message.Chat.ID
				responseText := fmt.Sprintf("Your chat id is: %d", chatID)
				msg := tgbotapi.NewMessage(chatID, responseText)
				if _, err := r.bot.Send(msg); err != nil {
					fmt.Printf("Failed to send /get_chat_id response: %v\n", err)
				}
			default:
				// Другие команды можно обрабатывать по мере необходимости.
			}
		}
	}
}
