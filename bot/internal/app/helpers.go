package app

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateInlineKeyBoard(buttons []string, prefix string) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for i := 0; i < len(buttons); i++ {
		rows[i] = make([]tgbotapi.InlineKeyboardButton, 0)
		rows[i] = append(rows[i], tgbotapi.NewInlineKeyboardButtonData(buttons[i], prefix+":"+buttons[i]))
	}

	return rows
}
