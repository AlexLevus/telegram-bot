package controllers

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

var membersCount int
var chatID int64

type BotController struct {
	bot *tele.Bot
}

func NewBotController(bot *tele.Bot) BotController {
	return BotController{bot}
}

func (bc *BotController) ShowHelpInfo(c tele.Context) error {
	return c.Send("Чтобы предложить фильм, используйте команду /suggest вместе с ссылкой на фильм из Кинопоиска")
}

func (bc *BotController) SuggestFilm(c tele.Context) error {
	linkToFilm := c.Message().Payload

	// ходим в api за информацией о фильме -> отправляем краткую выжимку
	// создаем голосование

	fmt.Printf(linkToFilm)

	// TODO переименовать
	chat := c.Message().Chat
	count, err := c.Bot().Len(chat)
	if err != nil {
		log.Fatal(err)
		return err
	}

	membersCount = count
	chatID = chat.ID

	// при создании чата привзяывать его id к чату

	poll := &tele.Poll{
		Type:     tele.PollRegular,
		Question: "Как насчет посмотреть \"Волк с Уолл Стрит\"?",
		Options:  []tele.PollOption{{Text: "Давай!"}, {Text: "Не хочу"}},
	}

	// r := &tele.ReplyMarkup{
	// 	 InlineKeyboard: [][]tele.InlineButton{{ { Text: "Давай!" } }},
	// }

	return c.Send(poll)
}

func (bc *BotController) HandlePollAnswer(c tele.Context) error {
	// если есть отрициательный ответ, то заканчиваем опрос и выводим сообщение

	poll := c.Poll()

	fmt.Printf("%+v\n", c.Chat())

	isPollEnded := poll.VoterCount == membersCount-1

	// TODO: если проголосовали все, кроме бота и создателя опроса
	if isPollEnded {
		isPollSuccessed := poll.Options[len(poll.Options)-1].VoterCount == 0
		chat, _ := bc.bot.ChatByID(chatID)

		if isPollSuccessed {
			_, err := bc.bot.Send(chat, "Отлично, все хотят посмотреть \"Волк с Уолл Стрит\"! Добавлю его в закладки")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := bc.bot.Send(chat, "Кому-то не понравился \"Волк с Уолл Стрит\". Предложи другой фильм")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}
