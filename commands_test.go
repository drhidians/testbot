package telegram

import "testing"

var UTF16SliceTests = []struct {
	S        string
	From, To int
	Slice    string
}{
	// English (ASCII)
	{"abc", 0, 0, ""},
	{"abc", 0, 1, "a"},
	{"abc", 0, 3, "abc"},
	{"abc", 1, 3, "bc"},
	{"abc", 3, 3, ""},
	{"abc", 0, -1, "abc"},
	// Russian
	{"абв", 0, 0, ""},
	{"абв", 0, 1, "а"},
	{"абв", 0, 3, "абв"},
	{"абв", 1, 3, "бв"},
	{"абв", 3, 3, ""},
	{"абв", 0, -1, "абв"},
}

func TestUTF16Slice(t *testing.T) {
	for _, tt := range UTF16SliceTests {
		if s := utf16Slice(tt.S, tt.From, tt.To); s != tt.Slice {
			t.Errorf("want %q, got %q", tt.Slice, s)
		}
	}
}

var SplitCommandTests = []struct {
	S, Command, Mention string
}{
	{"/command", "/command", ""},
	{"/command@bot", "/command", "bot"},
}

func TestSplitCommand(t *testing.T) {
	for _, tt := range SplitCommandTests {
		if c, m := splitCommand(tt.S); c != tt.Command || m != tt.Mention {
			t.Errorf("want (%q, %q), got (%q, %q)", tt.Command, tt.Mention, c, m)
		}
	}
}

var SplitArgsTests = []struct {
	S    string
	Args []string
}{
	{"", nil},
	{"", []string{}},
	{"    ", []string{}},
	{" a  b   ", []string{"a", "b"}},
}

func TestSplitArgs(t *testing.T) {
	for _, tt := range SplitArgsTests {
		if args := splitArgs(tt.S); !stringsEqual(args, tt.Args) {
			t.Errorf("want %s, got %s", tt.Args, args)
		}
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCommandRunRecovers(t *testing.T) {
	errc := make(chan error, 1)
	defer close(errc)
	c := command{Func: func(*Command, *Update) error { panic("test") }}
	c.Run(errc)
	err := <-errc
	if s := err.Error(); s != "test" {
		t.Errorf("error: want %q, got %q", "test", s)
	}
}
