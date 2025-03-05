package models

type User struct {
	TgID    int64  `json:"tg_id"`
	ChatID  int64  `json:"chat_id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	CityID  string `json:"city_id"`
	Region  string `json:"federal_subject,omitempty"`
	State   string `json:"state"`
	Sticker bool   `json:"sticker"`
}

func NewUser(tgID int64, chatID int64, name, state string) *User {
	return &User{
		TgID:   tgID,
		ChatID: chatID,
		Name:   name,
		State:  state,
	}
}

func (u *User) Update(city, cityID, state string, sticker bool, region string) {
	u.City = city
	u.CityID = cityID
	u.Region = region
	u.State = state
	u.Sticker = sticker
}
