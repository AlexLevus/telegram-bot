package routes

import (
	"github.com/AlexLevus/telegram-bot/internal/controllers"
	tele "gopkg.in/telebot.v3"
)

type BotRouteController struct {
	botController controllers.BotController
}

func NewBotRouteController(botController controllers.BotController) BotRouteController {
	return BotRouteController{botController}
}

func (rc *BotRouteController) BotRoute(bot *tele.Bot) {
	bot.Handle("/help", rc.botController.ShowHelpInfo)
	bot.Handle("/suggest", rc.botController.SuggestFilm)

	bot.Handle(tele.OnPoll, rc.botController.HandlePollAnswer)

	// в момент создания группы и добавления бота срабатывает OnAddedToGroup
	bot.Handle(tele.OnAddedToGroup, rc.botController.HandleAddedToGroup)

	// срабатывает в момент добавления пользователя
	bot.Handle(tele.OnUserJoined, rc.botController.HandleUserJoined)

	// срабатывает в момент удаления пользователя
	bot.Handle(tele.OnUserLeft, rc.botController.HandleUserLeft)
}