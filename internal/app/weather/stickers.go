package weather

import (
	"math/rand"
	"time"
	"weather-bot/internal/models"
)

var stickers = map[string][]string{
	"sunny": {"CAACAgIAAxkBAAEGDm1jRSTohtnDBEFuN2lsWPGmHhZMJAACxxkAAnm5qUqJP-wzIPZ3DyoE",
		"CAACAgUAAxkBAAEN6w5nwfJOtcT8oV0F1bwudaZgxnWTSAAC_AADcX38FGhiSj1aWThwNgQ",
		"CAACAgIAAxkBAAEN6xBnwfJ3-wlYeo1exStPZBFMJ7Ik3wACMFwAAv_mAUu0FDSk5qX5AAE2BA",
		"CAACAgIAAxkBAAEN6xZnwfLiJLaA2Llkz4YRaq33vSxHIwACRmkAAtJxgEpsLEc53u07GjYE",
		"CAACAgIAAxkBAAEN6xhnwfL7dXlgu4kyIwxwvFBqL0nnDAACklgAAveCiUqt_8YF-kdrfzYE",
		"CAACAgQAAxkBAAEN7W5nxD00CrIji49ShH6IXRYGfEzUQAACmgQAAngWfCkxfQzecaB-3TYE",
		"CAACAgIAAxkBAAEOH99n3CJG-_NWja7wz9v-1TUfHYhMBQACumAAAhe-sUgMBPtiJZkaOzYE"},
	"cloudy": {"CAACAgIAAxkBAAEGDl1jRR9BzBWweHPfV6fXrWW6uxLeWQACaRsAAnn7WEpAX5kHnzpXGioE",
		"CAACAgIAAxkBAAEN69hnwsuJanxoITNj54RVTFTqSksiWwACygUAAiMFDQABG-hgbSZuySA2BA",
		"CAACAgIAAxkBAAEOH9tn3CH0ntcSUTziOyB0A1aETxD7JgACCWoAAoYCKUjPeHp1U7SP3DYE",
		"CAACAgIAAxkBAAEOH91n3CIjHpjjrl2iWBSuXQHgBLA_TwACl2wAAlZsaEgym9OQp8cMjDYE"},
	"rain": {"CAACAgIAAxkBAAEGDmNjRSEnXTCCj4lXenCdHwHBNGvWRAACsBUAApWOQEgeOJ174sOmxSoE",
		"CAACAgIAAxkBAAEN6xRnwfKqPgRBCm9y-lnP5XzU-3FyagAC3E8AAjJ7OUhVF81S6WL2lDYE",
		"CAACAgIAAxkBAAEN6xpnwfMhjDIxoVfBMr7DyhbRXCghqQACzRoAAq8auUkq3NXKwblCFjYE",
		"CAACAgIAAxkBAAEN6x5nwfN7XLpOes1mO_Yr87qWKqL-UQACXx8AAhBYSEjOWqVjz1DnzDYE",
		"CAACAgIAAxkBAAEN69JnwsnS7fe2m5eFhvzYDMu-8eET2QACc0cAAj78OUioJLNfxKPfpzYE",
		"CAACAgIAAxkBAAEN7DFnwu_09vXcklW_VuQAAZAWBkpMF3QAAlwFAAIjBQ0AAScqEJ6OPBTxNgQ",
		"CAACAgIAAxkBAAEN7XRnxD2iYPJ_K-BOUEn40KQAAQRbXMUAAnEYAAI8DzBJ1DsPwUwgcIE2BA",
		"CAACAgIAAxkBAAEN7XhnxD38tCpNSfPk4tNxsobohZCPMQACPwoAAiD_YUoURq7SONcmGDYE",
		"CAACAgIAAxkBAAEN7XpnxD5-b27x6-hQQs4ZEVob2qa_LAACPDUAAqjo2EvVkfk_Cf8-uTYE"},
	"snow": {"CAACAgIAAxkBAAEGDmFjRSCyS4-BYboqYPMU-esJZIQsyAACrBQAAsQu2UssNI82XHFI0ioE",
		"CAACAgIAAxkBAAEGDmdjRSH0d85-ZzzPlngnv5nZX-k9kwACJBYAAts0yEujzSapZ8-iYyoE",
		"CAACAgIAAxkBAAEN6xJnwfKQdrQ6ECw-tKYk0_oa_wHGggACKUYAApMqOEi1syNCkZsuPTYE",
		"CAACAgIAAxkBAAEN6yJnwfQ1jKGTA2vGKguY9KhyLIH6TQAC1xkAArKIAAFLTPRaboiNs2w2BA",
		"CAACAgIAAxkBAAEN69BnwsmzSMHG8eVmoAYufPlAgLdAcQACPEkAAtOvOEiMLxxUIx9vfTYE",
		"CAACAgIAAxkBAAEN7Clnwu98i9NXM27V95-9ZwGVuayIFgACrAADU2LaKLrMyu2HOjpaNgQ",
		"CAACAgIAAxkBAAEN7Ctnwu-UbBuNFab-qdwpbBDaKlNq8gACsQADU2LaKI3r9VFlhdJvNgQ"},
	"windy": {"CAACAgIAAxkBAAEGDlFjRR3GKjAV6PxowHoXWJ8ZEC94qAACfRgAAsGAYUrNl1eKlj7diyoE",
		"CAACAgQAAxkBAAEN7XJnxD2Drm3JVSBEp_mS99KT6s_X_AACmAQAAngWfClFMNX_87UrHzYE"},
	"hot": {"CAACAgIAAxkBAAEGDp5jRSj0EujFDiRZTSFj8fUFaelkqgAC4RUAAgP1YEovHkwXccq2ESoE",
		"CAACAgIAAxkBAAEN6yBnwfQY4HlsTcxIgIY4AvYzV6GCDQAChBgAAvxNMUkp9tvuLH1wmjYE",
		"CAACAgIAAxkBAAEN68pnwsiWieIMJXPKjrK1GhSrLRbvgAACDyQAArtewEqY3xjZXL4skDYE",
		"CAACAgIAAxkBAAEN69RnwspedElk24cg-8O5LI0I141SVgACNRgAAuS2MUlCZC-mHrrHrDYE",
		"CAACAgIAAxkBAAEN69ZnwsqRmHVVdMXFs20DKEtopH21ewACORcAAnPNMEkfPXP8swRmlTYE",
		"CAACAgIAAxkBAAEN7C9nwu_RhUEzY1gC06JMtCdpAXlSnAAC7BkAApdj2ElUa-9zvJAbrDYE",
		"CAACAgQAAxkBAAEN7WxnxD0YZ0wfjx4gK6g_-VQEBXVt7gACNAQAAngWfCm258oLcK0ErTYE"},
	"cold": {"CAACAgIAAxkBAAEGDmVjRSGcWTq4xQg6VU78pndFP1mDMgACrxUAAqulyUtZecNIn_jluyoE",
		"CAACAgIAAxkBAAEN6xxnwfNCnsGVlg4CDoJISbzlqsGeLAAC-BoAAp8HuUkNd4LjiYDilzYE",
		"CAACAgIAAxkBAAEN685nwsk69KcQk21Hq7FWtB9bIZ0UlQACF0IAApuoOUiX5AKfgV51DDYE",
		"CAACAgIAAxkBAAEN7C1nwu-s010vnLBUJuMT3W35-aFkbwACbxQAApTa2EmyOApZ4L-0yjYE",
		"CAACAgIAAxkBAAEN7XZnxD3oyl2qTfmSeK7oWgTWpjIargAC9gUAAj6IGguoxaXu72HEkzYE"},
}

func getDominantConditionForSticker(weatherList map[int]int) string {
	var maxCount int
	var dominant string
	for condition, count := range weatherList {
		if count > maxCount {
			maxCount = count
			id := condition / 100
			switch {
			case condition == 800:
				dominant = "sunny"
			case id == 2, id == 3, id == 5:
				dominant = "rain"
			case id == 7, id == 8:
				dominant = "cloudy"
			case id == 6:
				dominant = "snow"
			}
		}
	}
	return dominant
}

func Sticker(forecast models.FullDayForecast) string {
	var groups [][]string

	weatherList := make(map[int]int)
	conditions := []int{
		forecast.Morning.ConditionId,
		forecast.Day.ConditionId,
		forecast.Evening.ConditionId,
		forecast.Night.ConditionId,
	}

	for _, conditionId := range conditions {
		if conditionId != 0 {
			weatherList[conditionId]++
		}
	}

	dominant := getDominantConditionForSticker(weatherList)
	groups = append(groups, stickers[dominant])
	if forecast.Day.WindSpeed > 4 {
		groups = append(groups, stickers["windy"])
	}
	if forecast.Day.Temperature > 20 {
		groups = append(groups, stickers["hot"])
	}
	if forecast.Day.Temperature < -3 {
		groups = append(groups, stickers["cold"])
	}

	if dominant == "" {
		return ""
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	stickers := groups[rand.Intn(len(groups))]
	return stickers[rand.Intn(len(stickers))]
}
