package bot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type capture struct {
	t                     *testing.T
	callCount             int
	helpHandlerCalled     bool
	sessionHandlerCalled  bool
	callbackHandlerCalled bool
	defaultHandlerCalled  bool
}

func (c *capture) helpHandler(b *Bot, u *UpdateResponse, args string) {
	c.helpHandlerCalled = true
	c.callCount += 1

	assert.Equal(c.t, 2, c.callCount)
	assert.Equal(c.t, "my options", args)
	assert.Equal(c.t, "/help@Test_Bot my options", u.Message.Text)
}

func (c *capture) sessionHandler(b *Bot, u *UpdateResponse, s SessionRecord) {
	c.sessionHandlerCalled = true
	c.callCount += 1

	assert.Equal(c.t, 1, c.callCount)
	assert.Equal(c.t, "this is my data", s.Data())
}

func (c *capture) callbackHandler(b *Bot, u *UpdateResponse) {
	c.callbackHandlerCalled = true
	c.callCount += 1

	assert.Equal(c.t, 1, c.callCount)
}

func (c *capture) defaultHandler(b *Bot, u *UpdateResponse, args string) {
	c.defaultHandlerCalled = true
	c.callCount += 1

	assert.Equal(c.t, 1, c.callCount)
	assert.Equal(c.t, "", args)
	assert.Equal(c.t, "non command message", u.Message.Text)
}

func newCapture(t *testing.T) *capture {
	return &capture{
		t: t,
	}
}

func runTest(t *testing.T, body string, checks func(c *capture)) {
	c := newCapture(t)
	s := newTestSession()

	s.SetSession(154355043, 145351026, 100, "this is my data")

	b := New("Test_Bot", "mysecrettoken")
	assert.NotNil(t, b)

	b.AddCommandHandler("help", c.helpHandler)
	b.AddSessionHandler(100, c.sessionHandler)
	b.SetDefaultHandler(c.defaultHandler)
	b.SetBeforeCommandCallback(c.callbackHandler)
	b.SetSession(s)

	req, _ := http.NewRequest("POST", "/bot10000000", strings.NewReader(body))
	b.HandleUpdate(req)

	checks(c)
}

func TestDefaultHandler(t *testing.T) {
	body := `{"update_id":797498290,"message":{"message_id":5265,"from":{"id":154355043,"first_name":"Tom","last_name":"Peters"},"date":1453565516,"chat":{"id":145351029,"type":"private","first_name":"John","last_name":"Doe"},"text":"non command message"}}`

	runTest(t, body, func(c *capture) {
		assert.False(t, c.callbackHandlerCalled)
		assert.False(t, c.helpHandlerCalled)
		assert.False(t, c.sessionHandlerCalled)
		assert.True(t, c.defaultHandlerCalled)
	})
}

func TestCommand(t *testing.T) {
	body := `{"update_id":258492060,"message":{"message_id":6261,"from":{"id":756606558,"first_name":"John","last_name":"Doe"},"date":1453514214,"chat":{"id":-974763016,"type":"group","title":"Test","first_name":""},"text":"/help@Test_Bot my options"}}`

	runTest(t, body, func(c *capture) {
		assert.True(t, c.callbackHandlerCalled)
		assert.True(t, c.helpHandlerCalled)
		assert.False(t, c.sessionHandlerCalled)
		assert.False(t, c.defaultHandlerCalled)
	})
}

func TestCommandForOtherBot(t *testing.T) {
	body := `{"update_id":258492060,"message":{"message_id":6261,"from":{"id":756606558,"first_name":"John","last_name":"Doe"},"date":1453514214,"chat":{"id":-974763016,"type":"group","title":"Test","first_name":""},"text":"/help@OtherTest_Bot my options"}}`

	runTest(t, body, func(c *capture) {
		// none of the handlers should be called
		assert.False(t, c.callbackHandlerCalled)
		assert.False(t, c.helpHandlerCalled)
		assert.False(t, c.sessionHandlerCalled)
		assert.False(t, c.defaultHandlerCalled)
	})
}

func TestSession(t *testing.T) {
	body := `{"update_id":797498290,"message":{"message_id":5265,"from":{"id":154355043,"first_name":"Tom","last_name":"Peters"},"date":1453565516,"chat":{"id":145351026,"type":"private","first_name":"John","last_name":"Doe"},"text":"Will this work?"}}`

	runTest(t, body, func(c *capture) {
		assert.False(t, c.callbackHandlerCalled)
		assert.False(t, c.helpHandlerCalled)
		assert.True(t, c.sessionHandlerCalled)
		assert.False(t, c.defaultHandlerCalled)
	})
}

func TestSessionForOtherChatID(t *testing.T) {
	body := `{"update_id":797498290,"message":{"message_id":5265,"from":{"id":154355043,"first_name":"Tom","last_name":"Peters"},"date":1453565516,"chat":{"id":145351027,"type":"private","first_name":"John","last_name":"Doe"},"text":"non command message"}}`

	runTest(t, body, func(c *capture) {
		assert.False(t, c.callbackHandlerCalled)
		assert.False(t, c.helpHandlerCalled)
		assert.False(t, c.sessionHandlerCalled)
		assert.True(t, c.defaultHandlerCalled)
	})
}

type testRoundTripper struct {
	request *http.Request
}

func (rt *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := strings.NewReader("")
	bodyCloser := ioutil.NopCloser(body)

	rt.request = r

	response := &http.Response{
		Status:     fmt.Sprintf("%d OK", http.StatusOK),
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       bodyCloser,
		Request:    r,
	}

	return response, nil
}

func TestPostSendMessage(t *testing.T) {
	transport := &testRoundTripper{}

	c := &http.Client{
		Transport: transport,
	}

	b := New("Test_Bot", "mysecrettoken")
	b.client = c

	msg := &SendMessage{
		ChatID:           1000,
		Text:             "This is my message",
		ReplyToMessageID: 500,
		ReplyMarkup: &ReplyMarkup{
			Keyboard:        [][]string{[]string{"One", "Two"}, []string{"Three", "Four"}},
			ResizeKeyboard:  false,
			OneTimeKeyboard: true,
			HideKeyboard:    false,
			ForceReply:      true,
			Selective:       false,
		},
	}

	err := b.PostSendMessage(msg)
	assert.NoError(t, err)

	defer transport.request.Body.Close()
	body, _ := ioutil.ReadAll(transport.request.Body)
	assert.Equal(t, `{"chat_id":1000,"text":"This is my message","reply_to_message_id":500,"reply_markup":{"keyboard":[["One","Two"],["Three","Four"]],"one_time_keyboard":true,"force_reply":true}}`, string(body[:len(body)-1]))
}

type testSessionRecord struct {
	authorID int
	chatID   int
	stateID  int
	data     string
}

func (r *testSessionRecord) AuthorID() int {
	return r.authorID
}

func (r *testSessionRecord) ChatID() int {
	return r.chatID
}

func (r *testSessionRecord) StateID() int {
	return r.stateID
}

func (r *testSessionRecord) Data() string {
	return r.data
}

type testSession struct {
	data map[string]*testSessionRecord
}

func newTestSession() *testSession {
	return &testSession{
		data: make(map[string]*testSessionRecord),
	}
}

func (s *testSession) key(authorID, chatID int) string {
	return fmt.Sprintf("%d.%d", authorID, chatID)
}

// SetSession should set a session for a user in a chat.
func (s *testSession) SetSession(authorID, chatID, stateID int, data string) error {
	s.data[s.key(authorID, chatID)] = &testSessionRecord{authorID, chatID, stateID, data}
	return nil
}

// DeleteSessionByAuthorIDAndChatID should delete a session for a user in a chat
func (s *testSession) DeleteSessionByAuthorIDAndChatID(authorID, chatID int) error {
	delete(s.data, s.key(authorID, chatID))
	return nil
}

// SessionByAuthorIDAndChatID should return a session for a user. If there is no session, but otherwise there was no error,
// (nil, nil) should be returned.
func (s *testSession) SessionByAuthorIDAndChatID(authorID, chatID int) (SessionRecord, error) {
	r, ok := s.data[s.key(authorID, chatID)]
	if !ok {
		return nil, nil
	}

	return r, nil
}
