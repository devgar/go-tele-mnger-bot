package main
import("os"; "os/exec"; "log"; "gopkg.in/telegram-bot-api.v4")
import("flag"; "fmt")

var bot *tgbotapi.BotAPI
var verbose bool
var scripts, user string
var token = "token"

func processUpdate(update tgbotapi.Update) {
  text, err := exec.Command("ls", "-s").Output()
  if err != nil { log.Panic(err) }
  
  msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(text))
  msg.ReplyToMessageID = update.Message.MessageID

  bot.Send(msg)
}

func initFlags() {
  flag.StringVar(&scripts, "scripts", "/etc/mngr/scripts", "scripts path")
  flag.StringVar(&token, "token", token, "Bot Api token")
  flag.StringVar(&user, "user", "", "Unique user allowed ID")
  flag.BoolVar(&verbose, "V", false, "Print additional infomation")
  flag.Parse()
  if verbose {
    log.Println("Arguments parsed:")
    fmt.Printf("  scripts '%s'\n", scripts)
    fmt.Printf("  token   '%s'\n", token)
    fmt.Printf("  user    '%s'\n", user)
  }
}

func initBot(key string) *tgbotapi.BotAPI {
  apiBot, err :=tgbotapi.NewBotAPI(key)
  if err != nil { log.Panic(err) }
  return apiBot
}

func main() {
  token = os.Getenv("MNGR_TOKEN")
  initFlags()
  bot = initBot(token)
  
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
