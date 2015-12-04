# telegram

Package telegram/bot handles interactions with the Telegram Bot API.

    * You can setup command handlers to respond to Telegram commands like /new
    * If you implement bot.Session, you can use an external database to
      maintain sessions for users between requests.

Here's a basic example:

    b := bot.New("Super_Bot", "TELEGRAM_API_KEY")

    b.AddCommandHandler("hello", func(b *bot.Bot, u *bot.UpdateResponse, args string) {
        msg := &bot.SendMessage{
            ChatID: u.ChatID(),
            Text:   fmt.Sprintf("Hello %s", args),
        }

        b.PostSendMessage(msg)
    })

    http.HandleFunc("/secretpath", func(w http.ResponseWriter, r *http.Request) {
        b.HandleUpdate(r)
    })

    log.Fatal(http.ListenAndServeTLS(":8444", "cert.pem", "key.pem", nil))
