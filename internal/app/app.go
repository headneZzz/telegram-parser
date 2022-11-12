package app

import (
	"context"
	"flag"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"os"
	"strconv"
	"telegram-parser/pkg/authentic"
)

// RunTelegram runs f callback with context and logger, panics on error.
func RunTelegram(f func(ctx context.Context, log *zap.Logger) error) {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()
	// No graceful shutdown.
	ctx := context.Background()
	if err := f(ctx, log); err != nil {
		log.Fatal("Run failed", zap.Error(err))
	}
}

func getChannelIdAndAccessHash(chatClasses []tg.ChatClass, channelName string, log *zap.Logger) (int64, int64) {
	var channelId int64
	var channelAccessHash int64
	for _, chatClass := range chatClasses {
		chatFull, ok := chatClass.AsFull()
		if !ok {
			return 0, 0
		}
		switch v := chatFull.(type) {
		case *tg.Chat:
		case *tg.Channel:
			channel := chatFull.(*tg.Channel)
			if channel.Title == channelName {
				log.Info(channel.Title)
				log.Info(strconv.FormatInt(channel.ID, 10))
				log.Info(strconv.FormatInt(channel.AccessHash, 10))
				channelId = channel.ID
				channelAccessHash = channel.AccessHash
			}
		default:
			panic(v)
		}
	}
	return channelId, channelAccessHash
}

func printMessages(messages *tg.MessagesChannelMessages, log *zap.Logger) {
	for _, message := range messages.Messages {
		switch v := message.(type) {
		case *tg.MessageEmpty:
		case *tg.Message:
			mes := message.(*tg.Message)
			log.Info(mes.Message)
		case *tg.MessageService:
		default:
			panic(v)
		}

	}
}

func getPhone() *string {
	envPhone := os.Getenv("PHONE")
	phone := flag.String("phone", envPhone, "phone number to authenticate")
	flag.Parse()
	return phone
}

func Run() {
	phone := getPhone()
	RunTelegram(func(ctx context.Context, log *zap.Logger) error {
		// Setting up authentication flow helper based on terminal authentic.
		flow := auth.NewFlow(
			authentic.TermAuth{UserPhone: *phone},
			auth.SendCodeOptions{},
		)
		client, err := telegram.ClientFromEnvironment(telegram.Options{
			Logger: log,
		})
		if err != nil {
			return err
		}
		return client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return err
			}
			var exceptIds []int64
			chats, err := client.API().MessagesGetAllChats(ctx, exceptIds)
			if err != nil {
				return err
			}
			chatClasses := chats.GetChats()
			channelName := "ТЕМКА В СХЕМКЕ"
			channelId, channelAccessHash := getChannelIdAndAccessHash(chatClasses, channelName, log)
			if channelId == 0 || channelAccessHash == 0 {
				return nil
			}
			request := tg.MessagesGetHistoryRequest{
				Peer: &tg.InputPeerChannel{
					ChannelID:  channelId,
					AccessHash: channelAccessHash,
				},
			}
			history, err := client.API().MessagesGetHistory(ctx, &request)
			if err != nil {
				return err
			}
			messages := history.(*tg.MessagesChannelMessages)
			printMessages(messages, log)
			return nil
		})
	})
}
