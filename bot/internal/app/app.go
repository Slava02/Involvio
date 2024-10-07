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
	"time"
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

	persistance := tm.NewLocalPersistence()

	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"get_user_data",
			persistance,
			tm.StateMap{
				"": {
					tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
						user := u.Update.SentFrom()

						file := tgbotapi.FilePath("/Users/slava/GolandProjects/Involvio/bot/internal/constants/img/Name.jpg")
						share := tgbotapi.NewPhoto(u.EffectiveChat().ID, file)
						bot.Send(share)

						msg := tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.NameMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton(fmt.Sprintf("%s", user.FirstName+" "+user.LastName)),
							),
						)
						bot.Send(msg)
						u.PersistenceContext.SetState("enter_gender")

						usr := &models.User{
							TelegID:  user.ID,
							UserName: user.UserName,
						}

						photo, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
							UserID: user.ID,
							Offset: 0,
							Limit:  1,
						})
						if err != nil {
							usr.Photo.FileID = ""
							log.Printf("Couldn't get profile photo: %s\n", err.Error())
						} else {
							if len(photo.Photos) != 0 {
								usr.Photo.FileID = photo.Photos[0][0].FileID
							} else {
								usr.Photo.FileID = ""
							}
						}
						u.PersistenceContext.PutDataValue(user.UserName, usr)
					}),
				},
				"enter_gender": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						msg := tgbotapi.NewMessage(u.EffectiveChat().ID, "👱‍ 👧 Укажи свой пол ")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						user.FullName = u.EffectiveMessage().Text

						msg = tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.GenderMsg,
						)
						msg.ReplyMarkup = constants.GenderCallback
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_city")
						u.PersistenceContext.PutDataValue(user.UserName, user)

					}),
				},
				"enter_city": {
					tm.NewCallbackQueryHandler(`^gender:(.+)`, nil, func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						gender := strings.Split(u.CallbackQuery.Data, ":")
						user.Gender = gender[1]

						msg := tgbotapi.NewMessage(
							u.CallbackQuery.Message.Chat.ID,
							constants.PhotoMsg,
						)
						if user.Photo.FileID != "" {
							msg.ReplyMarkup = constants.PhotoKeyBoard
						}
						bot.Send(msg)

						u.PersistenceContext.SetState("upload_photo")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"upload_photo": {
					tm.NewHandler(tm.Or(tm.HasPhoto(), tm.HasRegex("^"+constants.PhotoBtn)), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						if u.EffectiveMessage().Photo != nil {
							user.Photo.FileID = u.Message.Photo[0].FileID
						} else {
							photo, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
								UserID: user.TelegID,
								Offset: 0,
								Limit:  1,
							})
							if err != nil {
								log.Printf("Couldn't get profile photo: %s\n", err.Error())
							}
							user.Photo.FileID = photo.Photos[0][0].FileID
						}

						msg := tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							constants.CityMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_socials")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.EffectiveChat().ID,
							"Извините, не понял вас. Отправьте фотографию",
						))
					}),
				},
				"enter_socials": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.AgainBtn)), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						if user.City == "" {
							user.City = u.EffectiveMessage().Text
						}

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.SocialsMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_position")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"enter_position": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						_, err := url.ParseRequestURI(u.Message.Text)
						if err != nil {
							log.Println(err.Error())
							msg := tgbotapi.NewMessage(
								u.Message.Chat.ID,
								constants.WrongSocialsMsg,
							)
							msg.ReplyMarkup = constants.AgainKeyBoard
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
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"enter_interests": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						user.Position = u.Message.Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.InterestsMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_birthday")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"enter_birthday": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.AgainBtn)), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						user.Interests = u.Message.Text

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.BirthdayMsg,
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_goal")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"enter_goal": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						layout := "01/02/2006"
						t, err := time.Parse(layout, strings.Replace(u.Message.Text, ".", "/", -1))
						if err != nil {
							log.Printf("Не удалось получить время: %s\n", err.Error())
							msg := tgbotapi.NewMessage(
								u.Message.Chat.ID,
								constants.WrongTimeMsg,
							)
							msg.ReplyMarkup = constants.AgainKeyBoard
							bot.Send(msg)
							u.PersistenceContext.SetState("enter_birthday")
							return
						}
						user.Birthday = t.String()

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							constants.GoalMsg,
						)
						msg.ReplyMarkup = constants.GoalCallback
						bot.Send(msg)

						u.PersistenceContext.SetState("enter_group")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"enter_group": {
					tm.NewCallbackQueryHandler(`^goal:(.+)`, nil, func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

						goal := strings.Split(u.CallbackQuery.Data, ":")
						user.Goal = goal[1]

						msg := tgbotapi.NewMessage(
							u.CallbackQuery.Message.Chat.ID,
							constants.GroupCodeMsg,
						)
						msg.ReplyMarkup = constants.GroupKeyBoard
						bot.Send(msg)

						u.PersistenceContext.SetState("check_result")
						u.PersistenceContext.PutDataValue(user.UserName, user)
					}),
				},
				"check_result": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.GroupBtn)), func(u *tm.Update) {
						user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

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

						share := tgbotapi.NewPhoto(u.EffectiveChat().ID, tgbotapi.FileID(user.Photo.FileID))
						share.Caption = fmt.Sprintf(constants.ProfileMsg, user.FullName, user.City, user.UserName, user.Position, user.Interests)
						share.ReplyMarkup = constants.CheckResCallback
						bot.Send(share)

						u.PersistenceContext.SetState("final")
					}),
				},
				"final": {
					tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
						file := tgbotapi.FilePath("/Users/slava/GolandProjects/Involvio/bot/internal/constants/img/Reminder.jpg")
						share := tgbotapi.NewPhoto(u.EffectiveChat().ID, file)
						share.Caption = constants.FinalMsg
						share.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(share)

						u.PersistenceContext.SetState("")
					}),
				},
			},
			[]*tm.Handler{
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
				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.HelpMsg)
				msg.ReplyMarkup = constants.HelpCallback
				bot.Send(msg)
			},
		)).
		// TODO
		AddHandler(tm.NewCallbackQueryHandler(
			`^changeProfileData`,
			nil,
			func(u *tm.Update) {
				user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.NameMsg)
				//msg.ReplyMarkup = constants.HelpCallback
				bot.Send(msg)

				user.FullName = u.EffectiveMessage().Text

				u.PersistenceContext.PutDataValue(user.UserName, user)

				msg = tgbotapi.NewMessage(
					u.EffectiveChat().ID,
					fmt.Sprintf(constants.CheckProfileWithGroupMsg, u.Message.Text),
				)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
				bot.Send(msg)
			}),
		)

	for update := range updates {
		mux.Dispatch(bot, update)
	}

	return nil
}
