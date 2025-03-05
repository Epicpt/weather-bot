package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Telegram struct {
	Bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *Telegram {
	return &Telegram{Bot: bot}
}

func (t *Telegram) Message(chatID int64, text string, keyboard any) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	_, err := t.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка отправки сообщения")
	}
}

func (t *Telegram) Sticker(chatID int64, stickerID string) {
	msg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
	_, err := t.Bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка отправки стикера")
	}
}
