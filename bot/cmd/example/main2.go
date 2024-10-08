package example

import (
	"fmt"
	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const (
	StartShoppingMessage = "Start shopping"
	CheckoutMessage      = "Checkout"
)

type Cart map[string]bool

type Product struct {
	Sku   string
	Title string
	Image string
	Price int
}

var storage = map[string]Product{
	"p01": {
		Sku:   "a87xn",
		Title: "My awesome product",
		Image: "https://picsum.photos/id/252/600/400",
		Price: 89,
	},
	"p02": {
		Sku:   "a88xn",
		Title: "Another cool product",
		Image: "https://picsum.photos/id/253/600/400",
		Price: 1337,
	},
	"p03": {
		Sku:   "a89xn",
		Title: "The best product",
		Image: "https://picsum.photos/id/254/600/400",
		Price: 150000,
	},
}

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(StartShoppingMessage),
	),
)

var checkoutKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CheckoutMessage),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6845572405:AAHZ7DpYwZaA9Q-kv5fXQAOKogIo4XPodlw")
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {

				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Welcome to our shop!")
				msg.ReplyMarkup = startKeyboard

				bot.Send(msg)
			},
		)).
		AddHandler(tm.NewConversationHandler(
			"get_product_data",
			tm.NewLocalPersistence(),
			tm.StateMap{
				"": {
					tm.NewHandler(tm.And(tm.IsMessage(), tm.HasRegex("^"+StartShoppingMessage)), func(u *tm.Update) {
						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Provide a product ID. You can type /cancel at any time to cancel the process.",
						)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

						bot.Send(msg)
						u.PersistenceContext.SetState("enter_product_id")

						// Initialize user cart if not initialized yet
						if _, ok := u.PersistenceContext.GetData()["cart"]; !ok {
							u.PersistenceContext.PutDataValue("cart", make(Cart))
						}
					}),
				},
				"enter_product_id": {
					tm.NewHandler(tm.And(tm.IsMessage(), tm.HasRegex("^"+CheckoutMessage)), func(u *tm.Update) {
						cart := u.PersistenceContext.GetData()["cart"].(Cart)
						u.PersistenceContext.ClearData()

						var msg tgbotapi.MessageConfig

						if len(cart) > 0 {
							lines := []string{}
							for id := range cart {
								product := storage[id]
								lines = append(lines, fmt.Sprintf("- %s (%d $)", product.Title, product.Price))
							}

							msg = tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"Your has been recorded!\n\n"+strings.Join(lines, "\n")+"\n\nSee you again soon!",
							)
						} else {
							msg = tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"See you again soon!",
							)
						}

						msg.ReplyMarkup = startKeyboard
						bot.Send(msg)

						u.PersistenceContext.SetState("")
					}),
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						bot.Send(tgbotapi.NewChatAction(u.Message.Chat.ID, tgbotapi.ChatTyping))

						productID := u.Message.Text

						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")

						product, ok := storage[productID]
						if !ok {
							msg.Text = "Product not found"
							bot.Send(msg)
							return
						}

						file := tgbotapi.FileURL(product.Image)

						share := tgbotapi.NewPhoto(u.Message.Chat.ID, file)
						share.Caption = fmt.Sprintf("%s\n\nPrice: %d $\nSKU: %s", product.Title, product.Price, product.Sku)

						share.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Add to cart", "add:"+productID),
							),
						)

						bot.Send(share)

						instructions := tgbotapi.NewMessage(u.Message.Chat.ID, "Type another product ID to search, or click checkout to finish.")
						instructions.ReplyMarkup = checkoutKeyboard
						bot.Send(instructions)
					}),
					tm.NewHandler(tm.And(tm.Not(tm.IsCommandMessage("cancel")), tm.Not(tm.IsCallbackQuery())), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"ID only!",
						))
					}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")

					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "See you again soon!")
					msg.ReplyMarkup = startKeyboard

					bot.Send(msg)
				}),
				// During the active conversation these callback handler will be invoked
				// before the ones that are outside of this conversation.
				tm.NewCallbackQueryHandler(`^add:(.+)`, nil, func(u *tm.Update) {
					loadingMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Loading...", ""),
						),
					)

					bot.Send(tgbotapi.NewCallback(u.CallbackQuery.ID, "Refreshing..."))
					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						loadingMarkup,
					))
					bot.Send(tgbotapi.NewChatAction(u.CallbackQuery.Message.Chat.ID, tgbotapi.ChatTyping))

					productID := u.Context["matches"].([]string)[1]
					cart := u.PersistenceContext.GetData()["cart"].(Cart)
					cart[productID] = true
					u.PersistenceContext.PutDataValue("cart", cart)

					refreshMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Remove from cart", "remove:"+productID),
						),
					)

					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						refreshMarkup,
					))

					fmt.Printf("Cart: %v\n", cart)
				}),
				tm.NewCallbackQueryHandler(`^remove:(.+)`, nil, func(u *tm.Update) {
					loadingMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Loading...", ""),
						),
					)

					bot.Send(tgbotapi.NewCallback(u.CallbackQuery.ID, "Refreshing..."))
					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						loadingMarkup,
					))
					bot.Send(tgbotapi.NewChatAction(u.CallbackQuery.Message.Chat.ID, tgbotapi.ChatTyping))

					productID := u.Context["matches"].([]string)[1]
					cart := u.PersistenceContext.GetData()["cart"].(Cart)
					delete(cart, productID)
					u.PersistenceContext.PutDataValue("cart", cart)

					refreshMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Add to cart", "add:"+productID),
						),
					)

					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						refreshMarkup,
					))

					fmt.Printf("Cart: %v\n", cart)
				}),
			},
		)).
		AddHandler(tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
			callback := tgbotapi.NewCallback(u.CallbackQuery.ID, "Cannot modify cart at this time")
			if _, err := bot.Request(callback); err != nil {
				log.Print(err)
			}
			// Next:
			// Here you can handle any query callbacks even it from closed conversations.
		}))

	for update := range updates {
		mux.Dispatch(bot, update)
	}
}

//AddHandler(tm.NewConversationHandler(
//"change_profile_data_dialog",
//tm.NewLocalPersistence(), // we could also use `tm.NewFilePersistence("db.json")` or `&gormpersistence.GORMPersistence(db)` to keep data across bot restarts
//tm.StateMap{
//"": {
//tm.NewCallbackQueryHandler(
//`^changeProfileData`,
//nil,
//func(u *tm.Update) {
//user := storage.GetUser(u.EffectiveUser().ID)
//u.PersistenceContext.PutDataValue(user.UserName, user)
//
//share := tgbotapi.NewPhoto(u.EffectiveChat().ID, tgbotapi.FileID(user.Photo.FileID))
//share.Caption = fmt.Sprintf(constants.ProfileMsg+"\n"+constants.ChooseChangeDataOptMsg, user.FullName, user.City, user.UserName, user.Position, user.Interests)
//share.ReplyMarkup = constants.ChangeCallback
//bot.Send(share)
//
//u.PersistenceContext.SetState("enter_new_data")
//}),
//},
//"enter_new_data": {
//tm.NewCallbackQueryHandler(
//`^changeProfileData:(.+)`,
//nil,
//func(u *tm.Update) {
//var msg tgbotapi.MessageConfig
//switch strings.Split(u.CallbackData(), ":")[1] {
//case "fullName":
//u.PersistenceContext.PutDataValue("to_change", "fullName")
//msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.NameMsg)
//case "city":
//u.PersistenceContext.PutDataValue("to_change", "city")
//msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.CityMsg)
//case "position":
//u.PersistenceContext.PutDataValue("to_change", "position")
//msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.PositionMsg)
//case "interests":
//u.PersistenceContext.PutDataValue("to_change", "interests")
//msg = tgbotapi.NewMessage(u.EffectiveChat().ID, constants.InterestsMsg)
//default:
//log.Println("unknown options")
//}
//bot.Send(msg)
//
//u.PersistenceContext.SetState("change_data")
//}),
//tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))), func(u *tm.Update) {
//log.Println("2")
//bot.Send(tgbotapi.NewMessage(
//u.EffectiveChat().ID,
//"🤖 Извините, не понял вас.",
//))
//}),
//},
//"change_data": {
//tm.NewHandler(tm.HasText(), func(u *tm.Update) {
//user := u.PersistenceContext.GetData()[u.EffectiveUser().UserName].(*models.User)
//
//toChange := u.PersistenceContext.GetData()["to_change"].(string)
//
//switch toChange {
//case "fullName":
//user.FullName = u.Message.Text
//case "city":
//user.City = u.Message.Text
//case "position":
//user.Position = u.Message.Text
//case "interests":
//user.Interests = u.Message.Text
//default:
//log.Println("unknown option")
//}
//
//msg := tgbotapi.NewMessage(u.EffectiveChat().ID, "Успешно изменено!")
//bot.Send(msg)
//
//u.PersistenceContext.SetState("")
//}),
//tm.NewHandler(tm.Not(tm.Or(tm.IsCommandMessage("cancel"), tm.IsCommandMessage("start"))), func(u *tm.Update) {
//log.Println("1")
//bot.Send(tgbotapi.NewMessage(
//u.EffectiveChat().ID,
//"🤖 Извините, не понял вас.",
//))
//}),
//},
//},
//[]*tm.Handler{
//tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
//u.PersistenceContext.ClearData()
//u.PersistenceContext.SetState("")
//bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
//}),
//},
//))
