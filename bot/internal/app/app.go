package app

import (
	"fmt"
	"github.com/Slava02/Involvio/bot/config"
	"github.com/Slava02/Involvio/bot/internal/constants"
	"github.com/Slava02/Involvio/bot/internal/models"
	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/url"
	"strings"
)

const (
	serviceConfigKey = "Bot"
)

type Bot struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Bot {
	return &Bot{
		cfg: cfg,
	}
}

func (b *Bot) Start() error {
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"get_user_data",
			tm.NewLocalPersistence(),
			tm.StateMap{
				"": {
					tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
						file := tgbotapi.FilePath("/Users/slava/GolandProjects/Involvio/bot/internal/constants/img/Name.jpg")
						share := tgbotapi.NewPhoto(u.EffectiveChat().ID, file)
						bot.Send(share)

						msg := tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.NameMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						u.PersistenceContext.SetState("enter_gender")
						u.PersistenceContext.PutDataValue("user", new(models.User))
					}),
				},
				"enter_gender": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						user.FullName = u.EffectiveMessage().Text

						msg := tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.GenderMsg,
						)
						msg.ReplyMarkup = constants.GenderCallback
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_city")
						u.PersistenceContext.PutDataValue("user", user)

					}),
				},
				"enter_city": {
					tm.NewCallbackQueryHandler(`^gender:(.+)`, nil, func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						gender := strings.Split(u.CallbackQuery.Data, ":")
						user.Gender = gender[1]

						msg := tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.CityMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_socials")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_socials": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						user.City = u.EffectiveMessage().Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.SocialsMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_position")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_position": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						_, err := url.ParseRequestURI(u.Message.Text)
						if err != nil {
							log.Println(err.Error())
							msg := tgbotapi.NewMessage(
								u.Message.Chat.ID,
								constants.WrongSocialsMsg,
							)
							bot.Send(msg)
							u.PersistenceContext.SetState("enter_socials")
							return
						}
						user.Socials = u.Message.Text

						file := tgbotapi.FilePath("/Users/slava/GolandProjects/Involvio/bot/internal/constants/img/Position.jpg")
						share := tgbotapi.NewPhoto(u.Message.Chat.ID, file)
						bot.Send(share)

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.PositionMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_interests")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_interests": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						user.Interests = u.Message.Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.InterestsMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_birthday")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_birthday": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						user.Interests = u.Message.Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.BirthdayMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_goal")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_goal": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						// TODO: add validation
						user.Birthday = u.Message.Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.GoalMsg,
						)
						msg.ReplyMarkup = constants.GoalCallback
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_group")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"enter_group": {
					tm.NewCallbackQueryHandler(`^goal:(.+)`, nil, func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						goal := strings.Split(u.CallbackQuery.Data, ":")
						user.Goal = goal[1]

						msg := tgbotapi.NewMessage(
							u.CallbackQuery.Message.Chat.ID,
							constants.GroupCodeMsg,
						)
						msg.ReplyMarkup = constants.GroupKeyBoard
						bot.Send(msg)

						u.PersistenceContext.SetState("check_result")
						u.PersistenceContext.PutDataValue("user", user)
					}),
				},
				"check_result": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.GroupBtn)), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()["user"].(*models.User)

						var msg tgbotapi.MessageConfig
						if u.Message.Text == constants.GroupBtn {
							msg = tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								constants.CheckProfileMsg,
							)
						} else {
							msg = tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								fmt.Sprintf(constants.CheckProfileWithGroupMsg, u.Message.Text),
							)
						}
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						msg = tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							fmt.Sprintf(constants.ProfileMsg, user.FullName, user.City, user.Position, user.Interests),
						)
						msg.ReplyMarkup = constants.CheckResCallback
						bot.Send(msg)

						u.PersistenceContext.SetState("final")
					}),
				},
				"final": {
					tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
						file := tgbotapi.FilePath("/Users/slava/GolandProjects/Involvio/bot/internal/constants/img/Reminder.jpg")
						share := tgbotapi.NewPhoto(u.CallbackQuery.Message.Chat.ID, file)
						share.Caption = constants.FinalMsg
						share.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(share)
						u.PersistenceContext.SetState("")
					}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("start"), func(u *tm.Update) {
					log.Println("cleared context")
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")

					msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.StartMsg, u.Message.From.FirstName))
					msg.ReplyMarkup = constants.StartCallback

					bot.Send(msg)
				}),
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					log.Println("cleared context")
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")

					msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.StartMsg, u.Message.From.FirstName))
					msg.ReplyMarkup = constants.StartCallback

					bot.Send(msg)
				}),
				// During the active conversation these callback handler will be invoked
				// before the ones that are outside of this conversation.
			}),
		).
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.StartMsg, u.Message.From.FirstName))
				msg.ReplyMarkup = constants.StartCallback

				bot.Send(msg)
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("help"),
			func(u *tm.Update) {
				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, "")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
				bot.Send(msg)
				msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.HelpMsg)
				msg.ReplyMarkup = constants.HelpCallback
				bot.Send(msg)
			},
		))

	for update := range updates {
		mux.Dispatch(bot, update)
	}

	return nil
}
