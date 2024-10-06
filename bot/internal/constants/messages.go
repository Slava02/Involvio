package constants

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartMsg                 = "Привет, %s!👋\nЯ твой персональный ассистент 🤖\n\nКаждую неделю я буду предлагать тебе для встречи интересного человека, случайно выбранного среди других участников.\n\nДля старта ответь на несколько вопросов и прочитай короткую инструкциюю."
	NameMsg                  = "☕️ Как тебя зовут? Напиши имя и фамилию."
	CityMsg                  = "2/7.📍Из какого ты города? Напиши полное название города в ответ.\n✅ Санкт-Петербург, Ростов-на-Дону, Новосибирск\n❌ СПб, Ростов, Новосиб\n\nВы с партнером можете оказаться в разных городах – тогда вам подойдет онлайн формат встречи 👨‍💻\nЕсли вы из одного города – сможете встретиться вживую 🤝"
	SocialsMsg               = "👨‍💻 Пришли ссылку на профиль в соцсетях, который активно ведешь.\n💡 Ты можешь оставить ссылку на конкретный пост о себе или даже интервью — что-то, что поможет человеку заочно познакомиться с тобой до первой встречи. \n\n🔗 Оставляй сразу ссылку, по которой можно будет кликнуть. ❌ \"Мой никнейм в инсте – kukushe4ka\" – придётся набирать в поиске, это неудобно."
	PositionMsg              = "👨‍🔬 Кем ты работаешь и чем занимаешься?"
	InterestsMsg             = "👀 Какие у тебя есть рабочие и нерабочие интересы?\n\n💡 Напиши через запятую слова, за которые можно зацепиться и развернуть из этого интересный разговор :) Например, увлечения, города, названия книг."
	BirthdayMsg              = "📅 Напиши дату рождения в формате дд.мм.гггг."
	GoalMsg                  = "⚖️ Некоторые люди приходят на Random Coffee встречи, чтобы найти партнёров для будущих проектов и завести полезные контакты, условно назовём это \"пользой\". А кто-то приходит для расширения кругозора, новых эмоций и открытия чего-то нового, назовём это \"фан\". Какое описание больше подходит тебе?\n"
	GenderMsg                = "Укажи свой пол.\nВыбери подходящий вариант ниже 👇\n"
	CheckProfileMsg          = "Получилось! 🙌\n\nТеперь ты участник встреч Random Coffee ☕️\n\nВот так будет выглядеть твой профиль в сообщении, которое мы пришлем твоему собеседнику:\n⏬"
	CheckProfileWithGroupMsg = "Получилось! 🙌\n\nТеперь ты участник встреч Random Coffee ☕️\n\nУчастники будут подбираться только из группы с кодом: %s\n\nВот так будет выглядеть твой профиль в сообщении, которое мы пришлем твоему собеседнику:\n⏬"
	GroupCodeMsg             = "Если знаете, укажите код группы, тогда участники будут подбираться только в рамках этой группы"
	FinalMsg                 = "Хороших встреч! ☕️\n\n*  Свою пару для встречи ты будешь узнавать каждый понедельник — сообщение придет в этот чат\nНапиши участнику в Telegram, чтобы договориться о встрече или звонке.\nВремя и место вы выбираете сами.\n\n* В конце недели я спрошу: участвуешь ли ты на следующей неделе и как прошла твоя предыдущая встреча.\n"
	ProfileMsg               = "%s (%s)\nПрофиль: %s\n\nЧем занимается: %s\nЗацепки для начала разговора: %s\n\nЕсли нужно что-то поменять, поможет команда /help\n"
	HelpMsg                  = "Выбери подходящую опцию ниже. \nЕсли у тебя запрос посложнее, просто напиши его в ответ, и (может и не сразу, но обязательно) Random-coffee увидит его и ответит тебе. \n\nЧтобы перезапустить это меню, ты в любой момент можешь ввести /help"
	WrongSocialsMsg          = "Введите действительную ссылку"
)

// BUTTONS
const (
	StartBtn    = "Поехали!"
	GroupBtn    = "Не знаю"
	CheckResBtn = "Понятно, дальше"
	OkBtn       = "Ok"
)

// CALLBACKS
var (
	StartCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(StartBtn, "fill"),
		),
	)

	GenderCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мужчина", "gender:male"),
			tgbotapi.NewInlineKeyboardButtonData("Женщина", "gender:female"),
		),
	)

	GroupCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(GroupBtn, "none"),
		),
	)

	CheckResCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CheckResBtn, "next"),
		),
	)

	OkCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(OkBtn, "ok"),
		),
	)

	HelpCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть профиль", "checkProfile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поменять данные профиля", "changeProfileData"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выбрать еще сообщества", "changeGroups"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поставить бот на паузу", "pauseBot"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cнять бот с паузы", "resumeBot"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оценить встречу", "rateMeet"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заблокировать пользователя", "blockUser"),
		),
	)

	GoalCallback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Фан", "goal:fun"),
			tgbotapi.NewInlineKeyboardButtonData("Польза", "goal:benefits"),
			tgbotapi.NewInlineKeyboardButtonData("50/50", "goal:middle"),
		),
	)

	GroupKeyBoard = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GroupBtn)))
)
