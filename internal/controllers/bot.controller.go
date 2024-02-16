package controllers

import (
	"fmt"

	"github.com/AlexLevus/telegram-bot/internal/services"
	"github.com/AlexLevus/telegram-bot/internal/models"
	tele "gopkg.in/telebot.v3"
	"log"
)

var membersCount int
var chatID int64

type BotController struct {
	bot            *tele.Bot
	chatService    services.ChatService
}

func NewBotController(bot *tele.Bot, chatService services.ChatService) BotController {
	return BotController{bot, chatService}
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

	fmt.Printf("%+v\n", poll)

	return c.Send(poll)
}

func (bc *BotController) HandlePollAnswer(c tele.Context) error {
	// если есть отрициательный ответ, то заканчиваем опрос и выводим сообщение

	poll := c.Poll()

	isPollEnded := poll.VoterCount == membersCount-1

	// TODO: если проголосовали все, кроме бота и создателя опроса
	if isPollEnded {
		isPollSuccessed := poll.Options[len(poll.Options)-1].VoterCount == 0
		chat, _ := bc.bot.ChatByID(chatID)

		// bc.bot.StopPoll(poll.)

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

func (bc *BotController) HandleAddedToGroup(c tele.Context) error {
	chat := c.Chat()
	membersCount, err := c.Bot().Len(chat)

	if err != nil {
		log.Fatal(err)
	}

	addChatReq := models.SaveChatRequest{
		ChatId:       chat.ID,
		MembersCount: membersCount,
	}

	err = bc.chatService.SaveChat(&addChatReq)

	return err
}

func (bc *BotController) HandleMembersCountChange(c tele.Context) error {
	chat := c.Chat()
	membersCount, err := c.Bot().Len(chat)

	if err != nil {
		return err
	}

	return bc.chatService.UpdateChatMembersCount(chat.ID, membersCount)
}

func (bc *BotController) HandleUserJoined(c tele.Context) error {
	return bc.HandleMembersCountChange(c)
}

func (bc *BotController) HandleUserLeft(c tele.Context) error {
	return bc.HandleMembersCountChange(c)
}
