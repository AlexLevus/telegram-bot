package telegram_bot

import (
	"github.com/AlexLevus/telegram-bot/internal/repository"
	"log"
	"time"
	"os"
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func Start(rep *repository.Repository) {
	var membersCount int 
	var chatID int64

	counter, err := rep.GetCounter()
	fmt.Println(counter.Value)
	if err != nil {
		log.Fatalf("Error when get counter")
	}

	telegramApiToken, exists := os.LookupEnv("TELEGRAM_APITOKEN")
	if !exists {
		log.Fatal("Добавьте в файл .env api токен Телеграм")
	}


	pref := tele.Settings{
		Token:  telegramApiToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/help", func(c tele.Context) error {
		return c.Send("Чтобы предложить фильм, используйте команду /suggest вместе с ссылкой на фильм из Кинопоиска")
	})

	b.Handle("/suggest", func(c tele.Context) error {
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

		membersCount = count;
		chatID = chat.ID


		poll := &tele.Poll{
			Type:     tele.PollRegular,
			Question: "Как насчет посмотреть \"Волк с Уолл Стрит\"?",
			Options:  []tele.PollOption{{Text: "Давай!"}, {Text: "Не хочу"}},
		}

		// r := &tele.ReplyMarkup{
		// 	 InlineKeyboard: [][]tele.InlineButton{{ { Text: "Давай!" } }},
		// }

		return c.Send(poll)
	})


	b.Handle(tele.OnPoll, func(c tele.Context) error {
		// если есть отрициательный ответ, то заканчиваем опрос и выводим сообщение

		poll := c.Poll()

		fmt.Printf("%+v\n", poll)

		isPollEnded := poll.VoterCount == membersCount - 1

		// TODO: если проголосовали все, кроме бота и создателя опроса
		if isPollEnded {
			isPollSuccessed := poll.Options[len(poll.Options)-1].VoterCount == 0
			chat, _ := b.ChatByID(chatID)

			if isPollSuccessed {
				_, err = b.Send(chat, "Отлично, все хотят посмотреть \"Волк с Уолл Стрит\"! Добавлю его в закладки")
				return err
			} else {
				_, err = b.Send(chat, "Кому-то не понравился \"Волк с Уолл Стрит\". Предложи другой фильм")
				return err
			}
		}

		return errors.New("math: square root of negative number")
	})

	b.Start()
}