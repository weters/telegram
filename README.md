# telegram

Package telegram/bot handles interactions with the Telegram Bot API.

    * You can setup command handlers to respond to Telegram commands like /new
    * If you implement bot.Session, you can use an external database to
      maintain sessions for users between requests.

## Basic Example

Here's a basic example. We setup a command handler to listen for either "/hello <name>"
or "/hello@Super_Bot <name>". When the command is called, we will reply to user with "Hello <name>".

    b := bot.New("Super_Bot", "TELEGRAM_TOKEN")

    b.AddCommandHandler("hello", func(b *bot.Bot, u *bot.UpdateResponse, args string) {
        msg := &bot.SendMessage{
            ChatID: u.ChatID(),
            ReplyToMessageID: u.Message.ID,
            Text:   fmt.Sprintf("Hello %s", args),
        }

        b.PostSendMessage(msg)
    })

    http.HandleFunc("/secretpath", func(w http.ResponseWriter, r *http.Request) {
        b.HandleUpdate(r)
    })

    log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))

## Example With Session Handling

Here's a little more indepth one. We'll use the session mechanism so we can ask the user a question
and have the bot maintain state until the user responds.

    // Create two structs which will implement bot.Session and bot.SessionRecord respectively
    type Session struct{}
    type SessionRecord struct{}

    func (s *Session) SetSession(authorID, chatID, stateID int, data string) error {
        // should store data in a database
    }

    func (s *Session) DeleteSessionByAuthorIDAndChatID(authorID, chatID int) error {
        // should delete the data in a database
    }

    func (s *Session) SessionByAuthorIDAndChatID(authorID, chatID int) (SessionRecord, error) {
        // should return a session from the database
    }

    func (r *SessionRecord) AuthorID() int {
        // return author id
    }
    func (r *SessionRecord) ChatID() int {
        // return chat id
    }
    func (r *SessionRecord) StateID() int {
        // return state id
    }
    func (r *SessionRecord) Data() string {
        // return data
    }

    // keep track of the various sessions we'll want to store
    const (
        SAskFavoriteColor = iota
        SSomeOtherSession = iota
    )

    func main() {
        b := bot.New("Super_Bot", "TELEGRAM_TOKEN")

        b.SetSession(&Session{})

        // this command handler will respond to "/color"
        b.AddCommandHandler("color", func(b *bot.Bot, u *bot.UpdateResponse, args string) {
            msg := &bot.SendMessage{
                ChatID:           u.ChatID(),
                Text:             "What is your favorite color?",
                ReplyToMessageID: u.Message.ID,
                ReplyMarkup: &bot.ReplyMarkup{
                    Keyboard:        [][]string{[]string{"Red", "Blue"}, []string{"Green", "Yellow"}},
                    Selective:       true,
                    OneTimeKeyboard: true,
                },
            }

            db.SetSession(u.FromID(), u.ChatID(), SAskFavoriteColor, "")
            b.PostSendMessage(msg)
        })

        // this session handler will be called if the user didn't specify a command and they have
        // SAskFavoriteColor stored in the session.
        b.AddSessionHandler(SAskFavoriteColor, func(b *bot.Bot, u *bot.UpdateResponse, s bot.SessionRecord) {
            msg := &bot.SendMessage{
                ChatID: u.ChatID(),
                Text:   fmt.Sprintf("%s chose %s as their favorite color\n", u.Message.From.DisplayName(), u.Message.Text),
                ReplyMarkup: &bot.ReplyMarkup{
                    HideKeyboard: true,
                },
            }

            b.PostSendMessage(msg)
        })

        http.HandleFunc("/secretpath", func(w http.ResponseWriter, r *http.Request) {
            err := b.HandleUpdate(r)
            if err != nil {
                log.Printf("error: handle update had error: %s\n", err)
            }
        }

        log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
    }
