package main
import("os"; "os/exec"; "log"; "gopkg.in/telegram-bot-api.v4")

var bot *tgbotapi.BotAPI

func processUpdate(update tgbotapi.Update) {
  text, err := exec.Command("ls", "-s").Output()
  if err != nil { log.Panic(err) }
  
  msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(text))
  msg.ReplyToMessageID = update.Message.MessageID

  bot.Send(msg)
}

func initBot(key string) *tgbotapi.BotAPI {
  apiBot, err :=tgbotapi.NewBotAPI(key)
  if err != nil { log.Panic(err) }
  return apiBot
}

func main() {
  KEY := os.Getenv("MNGR_TOKEN")
  if len(os.Args) > 1 { KEY = os.Args[1] }
  bot = initBot(KEY)
  
  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates, err := bot.GetUpdatesChan(u)
  if err != nil { log.Panic(err) }

  for update := range updates {
    if update.Message == nil { continue }
    log.Printf("Message from %d", update.Message.From.ID)
    go processUpdate(update)
  }
}
