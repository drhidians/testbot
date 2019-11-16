package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testAPIResponse represents a alike apiResponse structure with Resutl key field.
type testAPIResponse struct {
	Response apiResponse
	Result   interface{}
}

// MarshalJSON implements json.Marshaler interface.
func (r *testAPIResponse) MarshalJSON() ([]byte, error) {
	var err error
	r.Response.Result, err = json.Marshal(r.Result)
	if err != nil {
		return nil, err
	}
	return json.Marshal(r.Response)
}

func TestBotDoIssuesValidHttpRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check HTTP method.
		if r.Method != "POST" {
			t.Errorf("method: want POST, got %s", r.Method)
		}
		// Check URL path.
		if r.URL.Path != "/token/method" {
			t.Fatalf("url: want %q, got %q", "/token/method", r.URL.Path)
		}
		// Check content.
		var m *Message
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			t.Fatal(err)
		}
		if *m.Text != "test" {
			t.Fatalf("text: want %q, got %q", "test", m.Text)
		}
		err := json.NewEncoder(w).Encode(&testAPIResponse{
			Response: apiResponse{OK: true},
			Result:   &User{Username: ref("bot")},
		})
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	ctx := context.Background()
	bot := newBot(ctx, "token", withURL(ts.URL+"/"))
	m := &Message{Text: ref("test")}
	var user User
	if err := bot.do(ctx, "method", m, &user); err != nil {
		t.Fatal(err)
	}
	if *user.Username != "bot" {
		t.Fatalf("username: want %q, got %q", "bot", user.Username)
	}
}

func ref(s string) *string {
	return &s
}

var updatesOptionsMarshalJSONTests = []struct {
	Opts updatesOptions
	JSON string
}{
	// Offset
	{updatesOptions{Offset: 1}, `{"offset":1}`},
	// Timeout must be encoded in seconds
	{updatesOptions{Timeout: time.Second}, `{"timeout":1}`},
	{updatesOptions{Timeout: time.Minute}, `{"timeout":60}`},
	// Limit
	{updatesOptions{Limit: 1}, `{"limit":1}`},
}

func TestUpdatesOptions_MarshalJSON(t *testing.T) {
	for _, tt := range updatesOptionsMarshalJSONTests {
		b, err := json.Marshal(&tt.Opts)
		if err != nil {
			t.Fatalf("%+v: %s", tt.Opts, err)
		}
		if s := string(b); s != tt.JSON {
			t.Fatalf("%+v: want %s, got %s", tt.Opts, tt.JSON, s)
		}
	}
}
