package telegram

import (
	"fmt"
	"strings"
	"unicode/utf16"
)

const mentionSign = "@"

// Command represents a command parsed from update's message.
// Args is a list of words right after the command in the message.
type Command struct {
	Name string
	Args []string
	From User
	Chat Chat
	Date int
}

// CommandFunc represents a function ran on every command.
// The function is ran in a separate goroutine.
type CommandFunc func(*Command, *Update) error

// Commands is the interface of a generic commands register/runner.
type Commands interface {
	Add(name string, fn CommandFunc)
	Run(*Update) (error, bool)
}

func NewCommands(username string) Commands {
	return &commands{username: username, m: map[string]CommandFunc{}}
}

type commands struct {
	username string
	m        map[string]CommandFunc
}

// Add adds the executor for the command. The executor will be called every time
// there will be its command in update.
func (c *commands) Add(name string, fn CommandFunc) {
	if !strings.HasPrefix(name, "/") {
		panic(fmt.Sprintf("telegram: command %q must start with /", name))
	}
	c.m[name] = fn
}

// Run creates commands from u and executes registered handlers for this commands.
// ok is false when no command is found in u. Check err whether handlers returned
// error while execution.
func (c *commands) Run(u *Update) (err error, ok bool) {
	if commands := c.parse(u); len(commands) != 0 {
		err = run(commands)
		ok = true
	}
	return
}

// parse returns a slice of known commands from u.
func (c *commands) parse(u *Update) []*command {
	if u == nil {
		return nil
	}
	if u.Message == nil {
		return nil
	}
	var cc []*command
	for _, e := range u.Message.Entities {
		if cmd := c.parseCommand(u.Message, e); cmd != nil {
			if fn, ok := c.m[cmd.Name]; ok {
				cc = append(cc, &command{Func: fn, Command: cmd, Update: u})
			}
		}
	}
	return cc
}

// parseCommand parses command from u update according to e message entity.
// If message entity is not a bot command or the command belongs to another bot
// then nil is returned.
//
// The API documentation garantees that Message.User will not be nil.
// So the pointer can be safely dereferenced.
func (c *commands) parseCommand(m *Message, e *MessageEntity) *Command {
	if !e.IsBotCommand() {
		return nil
	}
	text := *m.Text
	command := utf16Slice(text, e.Offset, e.Offset+e.Length)
	name, mention := splitCommand(command)
	// Commands to other bots are skipped.
	if mention != "" && mention != c.username {
		return nil
	}
	tail := utf16Slice(text, e.Offset+e.Length, -1)
	args := splitArgs(tail)
	return &Command{Name: name, Args: args, From: *m.From, Chat: m.Chat, Date: m.Date}
}

// splitCommand splits command from a message into command name and
// optional mention.
func splitCommand(s string) (name string, mention string) {
	pair := strings.SplitN(s, mentionSign, 2)
	name = pair[0]
	if len(pair) == 2 {
		mention = pair[1]
	}
	return
}

// splitArgs splits s into a list of arguments separated by spaces.
func splitArgs(s string) []string {
	s = strings.TrimSpace(s)
	return strings.Fields(s)
}

// run executes each command in a separate goroutine and returns multiError
// for all errors return by handlers.
func run(commands []*command) error {
	errc := make(chan error, len(commands))
	defer close(errc) // The channel will be drained before return.
	for i := range commands {
		go commands[i].Run(errc)
	}
	var err multiError
	for i := 0; i < cap(errc); i++ {
		err.Add(<-errc)
	}
	return err.Compact()
}

type command struct {
	Func    CommandFunc
	Command *Command
	Update  *Update
}

// Run executes the handler with the cmd and send error on errc.
func (c *command) Run(errc chan<- error) {
	var err error
	defer func() {
		// If Func panics then replace err with a recovered value - it will hold
		// the real error.
		if rerr := recover(); rerr != nil {
			err = fmt.Errorf("%s", rerr)
		}
		errc <- err
	}()
	err = c.Func(c.Command, c.Update)
}

// utf16Slice returns a substring of s from a UTF-16 code with from position
// to a code with to position.
// If to equals -1 then a slice will hold UTF-16 codes up to the last code.
//
// UTF-16 is used according to the documentation on offset and length fields.
// https://core.telegram.org/bots/api#messageentity
func utf16Slice(s string, from, to int) string {
	b := utf16.Encode([]rune(s))
	if to == -1 {
		to = len(b)
	}
	r := utf16.Decode(b[from:to])
	return string(r)
}

// multiError contains other errors.
type multiError struct {
	errs []error
}

// Add adds err into e.
func (e *multiError) Add(err error) {
	if err == nil {
		return
	}
	e.errs = append(e.errs, err)
}

// Compact returns e if e contains any error and nil otherwise.
func (e *multiError) Compact() error {
	if len(e.errs) == 0 {
		return nil
	}
	if len(e.errs) == 1 {
		return e.errs[0]
	}
	return e
}

// Error implements error interface.
func (e *multiError) Error() string {
	if len(e.errs) == 0 {
		return "<nil>"
	}
	s := make([]string, len(e.errs))
	for i := range e.errs {
		s[i] = e.errs[i].Error()
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ", "))
}
