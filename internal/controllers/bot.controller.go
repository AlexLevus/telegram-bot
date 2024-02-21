package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net/url"
	"os"
	"strings"

	"io"
	"net/http"

	"log"

	"github.com/AlexLevus/telegram-bot/internal/models"
	"github.com/AlexLevus/telegram-bot/internal/services"
	tele "gopkg.in/telebot.v3"
)

type BotController struct {
	bot         *tele.Bot
	chatService services.ChatService
	pollService services.PollService
	suggestionService services.SuggestionService
}

func NewBotController(bot *tele.Bot, chatService services.ChatService, pollService services.PollService, suggestionService services.SuggestionService) BotController {
	return BotController{bot, chatService, pollService, suggestionService}
}

func (bc *BotController) ShowHelpInfo(c tele.Context) error {
	return c.Send("Чтобы предложить фильм, используйте команду /suggest вместе с ссылкой на фильм из Кинопоиска")
}

type Response struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	Rating      struct {
		Kinopoisk float32 `json:"kp"`
	} `json:"rating"`
	Genres      []struct {
        Name        string    `json:"name"`
    } `json:"genres"`
	Poster struct {
		Url        string `json:"url"`
		PreviewUrl string `json:"previewUrl"`
	} `json:"poster"`
}

func (bc *BotController) GetFilm(filmID int) Response {
	url := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie/%d", filmID)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}

	kinopoiskApiToken, exists := os.LookupEnv("KINOPOISK_APITOKEN")
	if !exists {
		log.Fatal("Добавьте в файл .env api токен Телеграм")
	}

	req.Header.Set("X-API-KEY", kinopoiskApiToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf(resp.Status)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func (bc *BotController) SuggestFilm(c tele.Context) error {
	// ходим в api за информацией о фильме -> отправляем краткую выжимку
	// создаем голосование

	// получил id 
	// fmt.Printf("%+v\n", c.Message().Sender)

	filmUrl, err := url.ParseRequestURI(c.Message().Payload)
	if err != nil {
		return c.Send("Отправьте ссылку на фильм из Кинопоиска. Например, /suggest https://www.kinopoisk.ru/film/462682/", &tele.SendOptions{ DisableWebPagePreview: true })
	}

	stringFilmID := strings.Split(filmUrl.Path, "/")[2]
	filmID, _:= strconv.Atoi(stringFilmID)

	if (bc.suggestionService.IsSuggestionExists(c.Chat().ID, filmID)) {
		return c.Send("Этот фильм уже есть в ваших закладках")
	}


	film := bc.GetFilm(filmID)
	var genres []string

	for _, genre := range film.Genres {
		genres = append(genres, genre.Name)
    }

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Как насчет посмотреть \"%s\"?\n", film.Name))
	sb.WriteString(fmt.Sprintf("Год: %d. ", film.Year))
	sb.WriteString(fmt.Sprintf("Оценка: %.2f. ", film.Rating.Kinopoisk))
	sb.WriteString(fmt.Sprintf("Жанр: %s.", strings.Join(genres, ", ")))

	fmt.Println(sb.String())

	poll := &tele.Poll{
		Type:     tele.PollRegular,
		Question: sb.String(),
		Options:  []tele.PollOption{{Text: "Давай!"}, {Text: "Не хочу"}},
	}

	msg, err := c.Bot().Send(c.Recipient(), poll)
	if err != nil {
		return err
	}

	savePollReq := models.SavePollRequest{
		PollId:   msg.Poll.ID,
		ChatId:   msg.Chat.ID,
		MessageId: msg.ID,
		FilmName: film.Name,
		FilmId:   film.ID,
		FilmRating: film.Rating.Kinopoisk,
		FilmYear: film.Year,
		FilmUrl:  filmUrl.String(),
		Status: "processing",
	}

	return bc.pollService.SavePoll(&savePollReq)
}

func (bc *BotController) ShowSuggestions(c tele.Context) error {
	chat := c.Chat()
	
	suggestions, err := bc.suggestionService.FindSuggestionsByChatId(chat.ID)
	if err != nil {
		log.Fatal(err)
	}

	if len(suggestions) == 0 {
		return c.Send("В закладках пока нет ни одного фильма")
	}

	var sb strings.Builder
	sb.WriteString("Фильмы в закладках: \n\n")

    for idx, suggestion := range suggestions {
		sb.WriteString(fmt.Sprintf("%d. [%s](%s), %d, kp %.2f \n", idx + 1, suggestion.FilmName, suggestion.FilmUrl, suggestion.FilmYear, suggestion.FilmRating))
    }

	// fmt.Printf("%+v\n", suggestions[0])

	fmt.Println(sb.String())

	return c.Send(sb.String(), &tele.SendOptions{ ParseMode: tele.ModeMarkdown, DisableWebPagePreview: true })
}


type EditableMessage struct {
	messageID string
	chatID int64
}

func (m EditableMessage) MessageSig() (messageID string, chatID int64) {
	return m.messageID, m.chatID
}

func (bc *BotController) HandlePollAnswer(c tele.Context) error {
	poll := c.Poll()

	if poll.Closed {
		return nil
	}

	pollDB, _ := bc.pollService.FindPollById(c.Poll().ID)
	chatDB, _ := bc.chatService.FindChatById(pollDB.ChatId)

	fmt.Printf("%+v\n", poll.Options)

	hasPollNegativeAnswer := poll.Options[len(poll.Options)-1].VoterCount != 0
	chat, _ := bc.bot.ChatByID(pollDB.ChatId)

	pollEditableMessage := &EditableMessage{ messageID: strconv.Itoa(pollDB.MessageId), chatID: pollDB.ChatId }

	if hasPollNegativeAnswer {
		bc.pollService.UpdatePollStatus(pollDB.PollId, "failed")
		bc.bot.StopPoll(pollEditableMessage )
		_, err := bc.bot.Send(chat, fmt.Sprintf("Кому-то не понравился \"%s\". Предложи другой фильм", pollDB.FilmName))
		return err
	}

	isPollFinished := poll.VoterCount == chatDB.MembersCount-1

	if isPollFinished {
		bc.bot.StopPoll(pollEditableMessage)
		chat, _ := bc.bot.ChatByID(pollDB.ChatId)

		saveSuggestionReq := models.SaveSuggestionRequest{
			ChatId:   chat.ID,
			FilmId: pollDB.FilmId,
			FilmName: pollDB.FilmName,
			FilmUrl: pollDB.FilmUrl,
			FilmRating: pollDB.FilmRating,
			FilmYear: pollDB.FilmYear,
		}

		bc.suggestionService.SaveSuggestion(&saveSuggestionReq)
		bc.pollService.UpdatePollStatus(pollDB.PollId, "succeed")

		_, err := bc.bot.Send(chat, fmt.Sprintf("Отлично, все хотят посмотреть \"%s\"! Добавлю его в закладки", pollDB.FilmName))
		if err != nil {
			log.Fatal(err)
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

	saveChatReq := models.SaveChatRequest{
		ChatId:       chat.ID,
		MembersCount: membersCount,
	}

	err = bc.chatService.SaveChat(&saveChatReq)

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
