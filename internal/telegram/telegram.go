package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

type authService interface {
	AuthorizeTG(ctx context.Context, tgID int64, login, password string) error
	ValidateTGUser(ctx context.Context, tgID int64) (int, error)
}

type tgService interface {
	GetImageObjects(ctx context.Context, userId int) ([]io.Reader, error)
}

type Bot struct {
	tgService   tgService
	authService authService
	botAPI      *tgbotapi.BotAPI
	l           *logrus.Logger
}

const (
	reg      = "register"
	show     = "show"
	startCMD = "/start"
)

func NewBot(token string, l *logrus.Logger, tgService tgService, authService authService) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		botAPI:      bot,
		l:           l,
		tgService:   tgService,
		authService: authService,
	}, nil
}

func (b *Bot) StartBot() {
	bot := b.botAPI

	b.l.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.ProcessMessage(update.Message)

		} else if update.CallbackQuery != nil {
			b.l.Info(update.CallbackQuery.Data)
			chatId := update.CallbackQuery.Message.Chat.ID
			msgs := make([]tgbotapi.Chattable, 0)

			switch update.CallbackQuery.Data {
			case show:
				ctx := context.Background()
				userId, err := b.authService.ValidateTGUser(ctx, chatId)
				if err != nil {
					b.l.Error(err)
					msgs = append(msgs, tgbotapi.NewMessage(chatId, "Sign up!"))
					break
				}

				images, err := b.tgService.GetImageObjects(ctx, userId)
				if err != nil {
					b.l.Error(err)
				}

				for i, image := range images {
					byt, err := io.ReadAll(image)
					if err != nil {
						b.l.Error(err)
					}

					msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{
						Name:  strconv.Itoa(i) + ".jpg",
						Bytes: byt,
					})
					msgs = append(msgs, msg)
				}
			case reg:
				msgs = append(msgs, tgbotapi.NewMessage(chatId, "Enter your username and password.\nExample: test test"))
			}

			b.sendMsgs(msgs)
		}
	}
}

func (b *Bot) ProcessMessage(message *tgbotapi.Message) {
	var msg tgbotapi.MessageConfig
	chatId := message.Chat.ID
	switch message.Text {
	case startCMD:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Show images", show),
				tgbotapi.NewInlineKeyboardButtonData("Registration", reg),
			},
		)
		msg = tgbotapi.NewMessage(chatId, "Select the action")
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = keyboard
	default:
		var msgStr string
		s := strings.Split(message.Text, " ")
		if len(s) == 2 {
			err := b.authService.AuthorizeTG(context.Background(), message.From.ID, s[0], s[1])
			if err != nil {
				msgStr = "Incorrect login or password"
				b.l.Error(err)
			} else {
				msgStr = "You are registered!"
			}
		} else {
			msgStr = "Unknown command, enter:/start"
		}

		msg = tgbotapi.NewMessage(chatId, msgStr)
	}
	b.sendMsg(msg)
}

func (b *Bot) sendMsgs(msgs []tgbotapi.Chattable) {
	for _, msg := range msgs {
		b.sendMsg(msg)
	}
}

func (b *Bot) sendMsg(msg tgbotapi.Chattable) {
	_, err := b.botAPI.Send(msg)
	if err != nil {
		b.l.Error(err)
	}
}
