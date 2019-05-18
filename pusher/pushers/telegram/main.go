// Telegram Bot
// name: ServerMetricTestBot
// username: ServerMetricTest_bot
// token: 841563697:AAEDpNQBkNpFSUtae_ZgSRhxKzeJRntdrik
package telegram

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/darkfoxs96/servermetric/pusher"
)

type Telegram struct {
}

// Errors
var (
	ErrTokenDontString = fmt.Errorf("Telegram: field 'token' don't 'string' type")
	ErrDataDontString  = fmt.Errorf("Telegram: field 'data' don't 'string' type")
	ErrTokenNotFound   = fmt.Errorf("Telegram: field 'token' not found")
	ErrDataNotFound    = fmt.Errorf("Telegram: field 'data' not found")
)

func init() {
	if err := pusher.AppendPusher("telegram", &Telegram{}); err != nil {
		panic(err)
	}
}

func (t *Telegram) Init(config map[string]interface{}) (err error) {
	token := config["token"]
	if token == nil {
		return ErrTokenNotFound
	}

	data := config["data"]
	if data == nil {
		return ErrDataNotFound
	}

	tokenType := reflect.ValueOf(token)
	dataType := reflect.ValueOf(data)

	if tokenType.Kind() != reflect.String {
		return ErrTokenDontString
	}
	if dataType.Kind() != reflect.String {
		return ErrDataDontString
	}

	start(dataType.String())
	bot, err := tgbotapi.NewBotAPI(tokenType.String())
	if err != nil {
		return
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return
	}

	go watchAll(bot)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var messageText string

			switch update.Message.Text {
			case "/status":
				printWatchStatus(bot, update.Message.Chat.ID)
			case "/sub":
				subscribe(update.Message.Chat.ID)
				messageText = "You are subscribed now, Chat ID: " + strconv.Itoa(int(update.Message.Chat.ID)) // + "\n" + update.Message.Text
			case "/unsub":
				unsubscribe(update.Message.Chat.ID)
				messageText = "You are unsubscribed now, Chat ID: " + strconv.Itoa(int(update.Message.Chat.ID)) // + "\n" + update.Message.Text
			default:
				messageText = "Unknown command, use:\n" +
					"/sub to subscribe to notifications\n" +
					"/unsub to unsubscribe from notifications\n" +
					"/status to show current status"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}()

	return
}

func (t *Telegram) Push(msg string) (err error) {
	fmt.Println(msg)
	//msgQueue <- msg
	select {
	case msgQueue <- msg:

	default:
		fmt.Println("Warning! Can't write the message: ", msg)
	}

	return
}

var msgQueue = make(chan string, 200)

func watchAll(bot *tgbotapi.BotAPI) {
	for msg := range msgQueue {
		//fmt.Println("Got Telegram message: ", msg)
		config.mutex.Lock()
		sendNotification(bot, config.Subs, msg)
		config.mutex.Unlock()

	}
	//go watch(bot, &config.Subs, &config.mutex, migratorProcessName)
}

var (
	watchStatus     = make(map[string]bool)
	watchStatusMute sync.Mutex
)

func printWatchStatus(bot *tgbotapi.BotAPI, chatID int64) {
	var buf bytes.Buffer
	buf.WriteString("Current status:")
	watchStatusMute.Lock()
	for name, status := range watchStatus {
		buf.WriteString("\n" + name + ": ")
		if status {
			buf.WriteString("active")
		} else {
			buf.WriteString("inactive")
		}
	}
	watchStatusMute.Unlock()
	msg := tgbotapi.NewMessage(chatID, buf.String())
	bot.Send(msg)
}

func sendNotification(bot *tgbotapi.BotAPI, chatIDs []int64, messageText string) {
	msg := tgbotapi.NewMessage(0, messageText)
	for _, id := range chatIDs {
		msg.ChatID = id
		bot.Send(msg)
	}
}

func subscribe(id int64) {
	config.mutex.Lock()
	addID(&config.Subs, id)
	config.mutex.Unlock()

	config.save()
}

func addID(ids *[]int64, id int64) {
	for _, oid := range *ids {
		if oid == id {
			return
		}
	}
	*ids = append(*ids, id)
}

func unsubscribe(id int64) {
	config.mutex.Lock()
	removeID(&config.Subs, id)
	config.mutex.Unlock()

	config.save()
}

func removeID(ids *[]int64, id int64) {
	for i, oid := range *ids {
		if oid == id {
			*ids = append((*ids)[:i], (*ids)[i+1:]...)
			return
		}
	}
}
