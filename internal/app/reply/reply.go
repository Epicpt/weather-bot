package reply

type Sender interface {
	Message(chatID int64, text string, keyboard any)
	Sticker(chatID int64, stickerID string)
}

var sender Sender

func Init(s Sender) {
	sender = s
}

func Send() Sender {
	return sender
}
