package telegram

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestReplyKeyboardMarkup_MarshalJSON(t *testing.T) {
	m := &ReplyKeyboardMarkup{
		Keyboard: [][]*KeyboardButton{{{Text: "test"}}},
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	want := `"{\"keyboard\":[[{\"text\":\"test\"}]]}"`
	if s := string(b); s != want {
		t.Fatalf("json: want %q, got %q", want, s)
	}
}

func TestReplyKeyboardMarkup_UnmarshalJSON(t *testing.T) {
	var m ReplyKeyboardMarkup
	s := `"{\"keyboard\":[[{\"text\":\"test\"}]]}"`
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatal(err)
	}
	if s := m.Keyboard[0][0].Text; s != "test" {
		t.Fatalf("button: want %q, got %q", "test", s)
	}
}

var parseModeTests = []struct {
	Name  string
	Mode  ParseMode
	Bytes []byte
}{
	{"default", ModeDefault, []byte(`""`)},
	{"markdown", ModeMarkdown, []byte(`"Markdown"`)},
	{"html", ModeHTML, []byte(`"HTML"`)},
}

func TestParseMode_MarshalJSON(t *testing.T) {
	for _, tt := range parseModeTests {
		if b, err := json.Marshal(tt.Mode); err != nil {
			t.Fatalf("%s: %s", tt.Name, err)
		} else if !bytes.Equal(b, tt.Bytes) {
			t.Fatalf("%s: want %v, got %v", tt.Name, tt.Bytes, b)
		}
	}
}
