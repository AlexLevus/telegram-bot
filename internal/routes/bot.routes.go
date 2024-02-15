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
}