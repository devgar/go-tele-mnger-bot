package main

import s "strings"
import("os"; "os/exec"; "path"; "log"; "gopkg.in/telegram-bot-api.v4")
import("flag"; "fmt")

var bot *tgbotapi.BotAPI
var verbose, version bool
var scripts, token string
var user int

func versionEcho() {
  fmt.Println("v0.1-beta-2")
  os.Exit(0)
}

func processUpdate(update tgbotapi.Update) {
  if update.Message.From.ID != user {
    sendReply(update.Message, "Error: Unauthorized user")
    return
  }
  var txt = update.Message.Text
  if txt[0] != '_' {
    sendReply(update.Message, "Error: Unallowed")
    return
  }
  command := s.Split(txt[1:], " ")[0]

  if !scriptExists(command) {
    sendReply(update.Message, "Command not found")
    return
  }

  text, err := exec.Command(path.Join(scripts, command)).Output()
  if err != nil {
    sendReply(update.Message, "[!]Error on command execution")
  }
  sendReply(update.Message, "```\n" + string(text) + "\n```")
}

func sendReply(message *tgbotapi.Message, text string) {
  msg := tgbotapi.NewMessage(message.Chat.ID, text)
  msg.ReplyToMessageID = message.MessageID
  msg.ParseMode = "Markdown"
  bot.Send(msg)
}

func scriptExists(command string) bool {
  _, err := os.Stat(path.Join(scripts, command))
  return !os.IsNotExist(err)
}

func initFlags() {
  flag.StringVar(&scripts, "scripts", "/etc/mngr/scripts", "Scripts path")
  flag.StringVar(&token, "token", token, "Bot Api token")
  flag.IntVar(&user, "user", -1, "Unique user allowed ID")
  flag.BoolVar(&verbose, "V", false, "Print additional infomation")
  flag.BoolVar(&version, "v", false, "Show version")
  flag.Parse()
  if verbose {
    log.Println("Arguments parsed:")
    fmt.Printf("  scripts '%s'\n", scripts)
    fmt.Printf("  token   '%s'\n", token)
    fmt.Printf("  user    '%d'\n", user)
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
  if version {
    versionEcho()
  }
  bot = initBot(token)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates, err := bot.GetUpdatesChan(u)
  if err != nil { log.Panic(err) }

  for update := range updates {
    if update.Message == nil { continue }
    log.Printf("Message from %d (%s %s):\n%s\n",
      update.Message.From.ID,
      update.Message.From.FirstName,
      update.Message.From.LastName,
      update.Message.Text)
    go processUpdate(update)
  }
}
