package app

import (
	"fmt"
	"github.com/Slava02/Involvio/bot/config"
	"github.com/Slava02/Involvio/bot/internal/constants"
	"github.com/Slava02/Involvio/bot/internal/models"
	"github.com/Slava02/Involvio/bot/internal/repo"
	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
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
	//  INIT BOT API
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return err
	}

	//  INIT STORAGE
	storage := repo.New()
	storage.Data["s1av4_z"] = &models.User{
		TelegID:   670720852,
		FullName:  "Slava Zhuvaga",
		UserName:  "s1av4_z",
		Interests: "programming",
		Position:  "cook",
		Photo: models.Photo{
			FileID: "AgACAgEAAxUAAWcGlsVgWTNzWanDDpp7_QSRJNSQAALepzEbVGP6Jw8ThFSmJn5mAQADAgADYQADNgQ",
		},
		Groups: []string{"Общая"},
	}

	//  INIT UPDATES
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// INIT MUX
	mux := tm.NewMux().
		// START HANDLER
		//  TODO: add group parameters in initial query
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				const op = "START HANDLER"
				log := slog.With(
					slog.String("op", op),
					slog.AnyValue(u.Update),
				)
				log.Debug(op)

				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.StartMsg, u.Message.From.FirstName))
				msg.ReplyMarkup = constants.StartCallback

				bot.Send(msg)
			},
		)).
		//  GET USER DATA DIALOG
		AddHandler(tm.NewConversationHandler(
			"get_user_data",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  GET NAME
				"": {
					tm.NewCallbackQueryHandler(`fill`, nil,
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET NAME"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							//  SEND MESSAGE
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

							//  INIT USER
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
								slog.Debug("Couldn't get profile photo: %s\n", err.Error())
							} else {
								if len(photo.Photos) != 0 {
									usr.Photo.FileID = photo.Photos[0][0].FileID
								} else {
									usr.Photo.FileID = ""
								}
							}

							u.PersistenceContext.SetState("enter_gender")
							u.PersistenceContext.PutDataValue(user.UserName, usr)
						}),
				},
				//  GET GENDER
				"enter_gender": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET GENDER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

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

							u.PersistenceContext.SetState("enter_photo")
							u.PersistenceContext.PutDataValue(user.UserName, user)

						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET GENDER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  GET PHOTO
				"enter_photo": {
					tm.NewCallbackQueryHandler(`^gender:(.+)`, nil,
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET GET PHOTO"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

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

							u.PersistenceContext.SetState("enter_city")
							u.PersistenceContext.PutDataValue(user.UserName, user)
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET GET PHOTO"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  GET CITY
				"enter_city": {
					tm.NewHandler(tm.Or(tm.HasPhoto(), tm.HasRegex("^"+constants.PhotoBtn)),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET CITY"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							//  IF SENT PHOTO - USE IT
							if u.EffectiveMessage().Photo != nil {
								user.Photo.FileID = u.Message.Photo[0].FileID
							} else {
								photo, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
									UserID: user.TelegID,
									Offset: 0,
									Limit:  1,
								})
								if err != nil {
									slog.Debug("Couldn't get profile photo: %s\n", err.Error())
								}

								log.Info(fmt.Sprintf("PHOTO ID: %s\n", photo.Photos[0][0].FileID))
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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG: GET CITY"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас. Отправьте фотографию",
							))
						}),
				},
				//  ENTER SOCIALS
				"enter_socials": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.AgainBtn)),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER SOCIALS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER SOCIALS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  ENTER POSITION
				"enter_position": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER POSITION"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							//  TODO: validate links properly
							//  IF NOT VALID - BACK TO PREVIOUS STATE
							_, err := url.ParseRequestURI(u.Message.Text)
							if err != nil {
								slog.Debug(err.Error())
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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER POSITION"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  ENTER INTERESTS
				"enter_interests": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER INTERESTS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER INTERESTS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  ENTER BIRTHDAY
				"enter_birthday": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.AgainBtn)),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER BIRTHDAY"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER BIRTHDAY"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  ENTER GOAL
				"enter_goal": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER GOAL"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							//  VALIDATE DATE FORMAT
							layout := "01/02/2006"
							t, err := time.Parse(layout, strings.Replace(u.Message.Text, ".", "/", -1))
							//  IF NOT VALID - BACK TO PREVIOUS STATE
							if err != nil {
								slog.Debug("Не удалось получить время: %s\n", err.Error())
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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER GOAL"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  ENTER GROUP
				"enter_group": {
					tm.NewCallbackQueryHandler(`^goal:(.+)`, nil,
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER GROUP"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							//  TODO: add constants for saving goal
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
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  ENTER GROUP"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  CHECK PROFILE
				"check_result": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasRegex("^"+constants.GroupBtn)),
						func(u *tm.Update) {
							const op = "GET USER DATA DIALOG:  CHECK PROFILE"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							if u.Message.Text == constants.GroupBtn {
								storage.AddGroups(user.UserName, models.DefaultSpace)
							} else {
								newGroups := storage.AddGroups(user.UserName, u.EffectiveMessage().Text)
								if newGroups != "" {
									msg := tgbotapi.NewMessage(
										u.EffectiveChat().ID,
										fmt.Sprintf(constants.NewGroupsMsg, newGroups),
									)
									msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
									bot.Send(msg)
								}
							}

							log.Info("USER GROUPS: ", user.Groups)

							msg := tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								//  TODO: fix
								fmt.Sprintf(constants.CheckProfileWithGroupMsg, user.GetSpaces()),
							)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)

							share := tgbotapi.NewPhoto(u.EffectiveChat().ID, tgbotapi.FileID(user.Photo.FileID))
							share.Caption = fmt.Sprintf(constants.ProfileMsg, user.FullName, user.City, user.UserName, user.Position, user.Interests)
							bot.Send(share)

							storage.UpdateUser(user)

							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"),
					func(u *tm.Update) {
						slog.Debug("cleared context")
						u.PersistenceContext.ClearData()
						u.PersistenceContext.SetState("")

						bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Отменено"))
					}),
			}),
		).
		//  HELP
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("help"),
			func(u *tm.Update) {
				const op = "HELP HANDLER"
				log := slog.With(
					slog.String("op", op),
					slog.AnyValue(u.Update),
				)
				log.Debug(op)

				msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.HelpMsg)
				msg.ReplyMarkup = constants.HelpCallback
				bot.Send(msg)
			},
		)).
		//  HELP: CHECK_PROFILE
		AddHandler(tm.NewCallbackQueryHandler("checkProfile", nil,
			func(u *tm.Update) {
				const op = "HELP: CHECK_PROFILE HANDLER"
				log := slog.With(
					slog.String("op", op),
					slog.AnyValue(u.Update),
				)
				log.Debug(op)
				user := storage.GetUser(u.EffectiveUser().UserName)

				share := tgbotapi.NewPhoto(u.EffectiveChat().ID, tgbotapi.FileID(user.Photo.FileID))
				share.Caption = fmt.Sprintf(constants.ProfileMsg, user.FullName, user.City, user.UserName, user.Position, user.Interests)
				bot.Send(share)
			},
		)).
		//  HELP: CHANGE_PROFILE_DIALOG
		AddHandler(tm.NewConversationHandler(
			"change_profile_data_dialog",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  CHOOSE CHANGE OPTION
				"": {
					tm.NewCallbackQueryHandler(`^changeProfileData`, nil,
						func(u *tm.Update) {
							const op = "HELP: CHANGE_PROFILE_DIALOG: CHOOSE CHANGE OPTION"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := storage.GetUser(u.EffectiveUser().UserName)
							u.PersistenceContext.PutDataValue(user.UserName, user)

							share := tgbotapi.NewPhoto(u.EffectiveChat().ID, tgbotapi.FileID(user.Photo.FileID))
							share.Caption = fmt.Sprintf(constants.ProfileMsg+"\n"+constants.ChooseChangeDataOptMsg, user.FullName, user.City, user.UserName, user.Position, user.Interests)
							share.ReplyMarkup = constants.ChangeCallback
							bot.Send(share)

							u.PersistenceContext.SetState("enter_new_data")
						}),
				},
				// ENTER NEW DATA
				"enter_new_data": {
					tm.NewCallbackQueryHandler(`^changeProfileData:(.+)`, nil,
						func(u *tm.Update) {
							const op = "HELP: CHANGE_PROFILE_DIALOG: ENTER NEW DATA"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := storage.GetUser(u.EffectiveUser().UserName)

							var msg tgbotapi.MessageConfig
							//  DEPENDING ON USER CHOICE - WE SEND DIFFERENT MESSAGES
							//  AND SAVE INFO ABOUT WHAT HAS TO BE CHANGED
							//  TODO: move out key-gen
							key := "to_change" + user.UserName
							switch strings.Split(u.CallbackData(), ":")[1] {
							case "fullName":
								u.PersistenceContext.PutDataValue(key, "fullName")
								msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.NameMsg)
							case "city":
								u.PersistenceContext.PutDataValue(key, "city")
								msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.CityMsg)
							case "position":
								u.PersistenceContext.PutDataValue(key, "position")
								msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.PositionMsg)
							case "interests":
								u.PersistenceContext.PutDataValue(key, "interests")
								msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.InterestsMsg)
							case "photo":
								u.PersistenceContext.PutDataValue(key, "photo")
								msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.PhotoMsg)
							default:
								slog.Debug("unknown options")
							}
							bot.Send(msg)

							u.PersistenceContext.SetState("change_data")
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: CHANGE_PROFILE_DIALOG: ENTER NEW DATA"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
				//  CHANGE DATA
				"change_data": {
					tm.NewHandler(tm.Or(tm.HasText(), tm.HasPhoto()),
						func(u *tm.Update) {
							const op = "HELP: CHANGE_PROFILE_DIALOG: CHANGE DATA"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)
							key := "to_change" + user.UserName
							//  GET INFO ABOUT WHAT DO WE NEED TO CHANGE
							toChange := u.PersistenceContext.GetData()[key].(string)

							switch toChange {
							case "fullName":
								user.FullName = u.Message.Text
							case "city":
								user.City = u.Message.Text
							case "position":
								user.Position = u.Message.Text
							case "interests":
								user.Interests = u.Message.Text
							case "photo":
								slog.Debug("PHOTO: %s\n", u.Message.Photo[0].FileID)
								user.Photo.FileID = u.Message.Photo[0].FileID
							default:
								slog.Debug("unknown option")
							}

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.ChangedMsg)
							bot.Send(msg)

							//  UPDATE USER INFO
							storage.UpdateUser(user)

							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: CHANGE_PROFILE_DIALOG: CHANGE DATA"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Отменено"))
				}),
			},
		)).
		//  HELP: CHANGE GROUPS DIALOG
		AddHandler(tm.NewConversationHandler(
			"change_groups_dialog",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  ENTER GROUPS
				"": {
					tm.NewCallbackQueryHandler(`^changeGroups`, nil,
						func(u *tm.Update) {
							const op = "HELP: CHANGE GROUPS DIALOG: ENTER GROUPS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := storage.GetUser(u.EffectiveUser().UserName)
							u.PersistenceContext.PutDataValue(user.UserName, user)

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.GroupCodeMsg)
							bot.Send(msg)

							u.PersistenceContext.SetState("update_groups")
						}),
				},
				//  UPDATE GROUPS
				"update_groups": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "HELP: CHANGE GROUPS DIALOG: UPDATE GROUPS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							newGroups := storage.AddGroups(user.UserName, u.EffectiveMessage().Text)
							if newGroups != "" {
								msg := tgbotapi.NewMessage(
									u.EffectiveChat().ID,
									fmt.Sprintf(constants.NewGroupsMsg, newGroups),
								)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
								bot.Send(msg)
							}

							log.Info("USER GROUPS: ", user.Groups)

							msg := tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								fmt.Sprintf(constants.ChangedGroupMsg, user.GetSpaces()))
							bot.Send(msg)

							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: CHANGE GROUPS DIALOG: UPDATE GROUPS"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Отменено"))
				}),
			},
		)).
		//  HELP: RESUME BOT
		AddHandler(tm.NewCallbackQueryHandler("resumeBot", nil,
			func(u *tm.Update) {
				const op = "HELP: CHANGE GROUPS DIALOG: UPDATE GROUPS"
				log := slog.With(
					slog.String("op", op),
					slog.AnyValue(u.Update),
				)
				log.Debug(op)

				// TODO: change state
				bot.Send(tgbotapi.NewMessage(u.EffectiveChat().ID, constants.BotResumedMsg))
			},
		)).
		//  HELP: PAUSE BOT DIALOG
		AddHandler(tm.NewConversationHandler(
			"pause_bot_dialog",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  CHOOSE PAUSE OPTION
				"": {
					tm.NewCallbackQueryHandler(`^pauseBot`, nil,
						func(u *tm.Update) {
							const op = "HELP: PAUSE BOT DIALOG: CHOOSE PAUSE OPTION"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							user := storage.GetUser(u.EffectiveUser().UserName)
							u.PersistenceContext.PutDataValue(user.UserName, user)

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.PauseBotOprionsMsg)
							msg.ReplyMarkup = constants.PauseOptionsCallBack
							bot.Send(msg)

							u.PersistenceContext.SetState("enter_pause")
						}),
				},
				//  PAUSE BOT
				"enter_pause": {
					tm.NewCallbackQueryHandler(`^pauseBot:(.+)`, nil,
						func(u *tm.Update) {
							const op = "HELP: PAUSE BOT DIALOG: PAUSE BOT"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)
							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.BotPausedMsg)
							bot.Send(msg)
							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: PAUSE BOT DIALOG: PAUSE BOT"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас.",
							))
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
				}),
			},
		)).
		//  HELP: RATE MEET DIALOG
		AddHandler(tm.NewConversationHandler(
			"rate_meet_dialog",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  ENTER USER
				"": {
					tm.NewCallbackQueryHandler(`^rateMeet`, nil,
						func(u *tm.Update) {
							const op = "HELP: RATE MEET DIALOG: ENTER USER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := storage.GetUser(u.EffectiveUser().UserName)

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.ChooseUserToRateMsg)
							msg.ReplyMarkup = constants.RateLastMeetKeyBoard
							bot.Send(msg)

							u.PersistenceContext.SetState("rate_user")
							u.PersistenceContext.PutDataValue(user.UserName, user)
						}),
				},
				//  RATE USER
				"rate_user": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "HELP: RATE MEET DIALOG: RATE USER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)
							//  TODO: move key generation out
							key := "user_to_rate" + user.UserName

							var msg tgbotapi.MessageConfig
							//  TODO: add check if user exists
							//  IF USER PICKED TO RATE LAST MEET
							if u.EffectiveMessage().Text == constants.RateLastBtn {
								//  TODO: check if user had meetings
								if x := (time.Now().Second()) % 2; x == 0 {
									//  THERE WERE NO MEETINGS YET
									msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.NoMeetsYetMsg)
									msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
									bot.Send(msg)

									u.PersistenceContext.ClearData()
									u.PersistenceContext.SetState("")
									return
								} else {
									//  TODO: get last meeting
									u.PersistenceContext.PutDataValue(key, "@s1av4")
								}
							} else {
								//  TODO: validate nickname
								u.PersistenceContext.PutDataValue(key, u.EffectiveMessage().Text)
							}

							msg = tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.RateUserMsg, u.PersistenceContext.GetData()[key].(string)))
							msg.ReplyMarkup = constants.RateMeetCallback
							bot.Send(msg)

							u.PersistenceContext.SetState("save_user_rate")
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: RATE MEET DIALOG: RATE USER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас. Введите пользователя в формате @ник",
							))
						}),
				},
				//  SAVE RATE
				"save_user_rate": {
					tm.NewCallbackQueryHandler(`^rateMeet:(.+)`, nil,
						func(u *tm.Update) {
							const op = "HELP: RATE MEET DIALOG: SAVE RATE"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)

							key := "user_to_rate" + user.UserName
							user_to_rate := u.PersistenceContext.GetData()[key].(string)

							//TODO: save rate

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.UserRatedMsg, user_to_rate))
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)

							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: RATE MEET DIALOG: SAVE RATE"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас. Введите пользователя в формате @ник",
							))
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
				}),
			},
		)).
		//  HELP: BLOCK USER DIALOG
		AddHandler(tm.NewConversationHandler(
			"block_user_dialog",
			tm.NewLocalPersistence(),
			tm.StateMap{
				//  ENTER USER TO BLOCK
				"": {
					tm.NewCallbackQueryHandler(`^blockUser`, nil,
						func(u *tm.Update) {
							const op = "HELP: BLOCK USER DIALOG: ENTER USER TO BLOCK"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							user := storage.GetUser(u.EffectiveUser().UserName)
							u.PersistenceContext.PutDataValue(user.UserName, user)

							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, constants.BlockUserMsg)
							bot.Send(msg)

							u.PersistenceContext.SetState("block_user")
						}),
				},
				//  BLOCK USER
				"block_user": {
					tm.NewHandler(tm.HasText(),
						func(u *tm.Update) {
							const op = "HELP: BLOCK USER DIALOG: BLOCK USER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							//  TODO: check if user exists
							msg := tgbotapi.NewMessage(u.EffectiveChat().ID, fmt.Sprintf(constants.UserBlockedMsg, u.EffectiveMessage().Text))
							bot.Send(msg)

							//  TODO: save info

							u.PersistenceContext.SetState("")
							u.PersistenceContext.ClearData()
						}),
					tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))),
						func(u *tm.Update) {
							const op = "HELP: BLOCK USER DIALOG: BLOCK USER"
							log := slog.With(
								slog.String("op", op),
								slog.AnyValue(u.Update),
							)
							log.Debug(op)

							bot.Send(tgbotapi.NewMessage(
								u.EffectiveChat().ID,
								"🤖 Извините, не понял вас. Введите пользователя в формате @ник",
							))
						}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Отменено"))
				}),
			},
		))

	for update := range updates {
		mux.Dispatch(bot, update)
	}

	return nil
}
