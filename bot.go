package main

import s "strings"
import("os"; "os/exec"; "path"; "log"; "gopkg.in/telegram-bot-api.v4")


var bot *tgbotapi.BotAPI

func processUpdate(update tgbotapi.Update) {
  var txt = update.Message.Text
  if txt[0] != '_' {
    sendReply(update.Message, "Error: Unallowed")
  } else {
    command := s.Split(txt[1:], " ")[0] + ".sh"
    log.Printf("Command to run: %s", command)
    text, err := exec.Command("ls", "-s").Output()
    if err != nil { log.Panic(err) }
    sendReply(update.Message, string(text))
  }
}

func sendReply(message *tgbotapi.Message, text string) {
  msg := tgbotapi.NewMessage(message.Chat.ID, text)
  msg.ReplyToMessageID = message.MessageID
  
  bot.Send(msg)
}

func scriptExists(command string) bool {
  _, err := os.Stat(path.Join("/usr", command))
  return os.IsNotExist(err)
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
    log.Printf("Message from %d (%s %s):\n%s\n", 
      update.Message.From.ID,
      update.Message.From.FirstName,
      update.Message.From.LastName, 
      update.Message.Text)
    go processUpdate(update)
  }
}
