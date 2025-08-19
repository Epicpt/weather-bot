package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *Telegram {
	return &Telegram{Bot: bot}
}

func (t *Telegram) Message(chatID int64, text string, keyboard any) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (t *Telegram) Sticker(chatID int64, stickerID string) error {
	msg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
