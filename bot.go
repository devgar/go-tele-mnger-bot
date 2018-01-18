package main

import s "strings"
import("os"; "os/exec"; "path"; "log"; "gopkg.in/telegram-bot-api.v4")
import("flag"; "fmt")

var bot *tgbotapi.BotAPI
var verbose bool
var scripts, user string
var token = "token"

func processUpdate(update tgbotapi.Update) {
  var txt = update.Message.Text
  if txt[0] != '_' {
    sendReply(update.Message, "Error: Unallowed")
    return 
  }
  command := s.Split(txt[1:], " ")[0] + ".sh"
  
  if !scriptExists(command) {
    sendReply(update.Message, "Command not found")
    return 
  }
  
  text, err := exec.Command(command).Output()
  if err != nil {
    sendReply(update.Message, "[!]Error on command execution")
  }
  sendReply(update.Message, string(text))
}

func sendReply(message *tgbotapi.Message, text string) {
  msg := tgbotapi.NewMessage(message.Chat.ID, text)
  msg.ReplyToMessageID = message.MessageID
  
  bot.Send(msg)
}

func scriptExists(command string) bool {
  _, err := os.Stat(path.Join(scripts, command))
  return !os.IsNotExist(err)
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
    log.Printf("Message from %d (%s %s):\n%s\n", 
      update.Message.From.ID,
      update.Message.From.FirstName,
      update.Message.From.LastName, 
      update.Message.Text)
    go processUpdate(update)
  }
}
