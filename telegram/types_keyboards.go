package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "encoding/json"

type aliasReplyKeyboardMarkup ReplyKeyboardMarkup
type aliasReplyKeyboardRemove ReplyKeyboardRemove
type aliasInlineKeyboardMarkup InlineKeyboardMarkup
type aliasForceReply ForceReply

// MarshalJSON implements json.Marshaler interface.
func (m *ReplyKeyboardMarkup) MarshalJSON() ([]byte, error) {
	a := (*aliasReplyKeyboardMarkup)(m)
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (m *ReplyKeyboardMarkup) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var a aliasReplyKeyboardMarkup
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	}
	*m = (ReplyKeyboardMarkup)(a)
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (m *ReplyKeyboardRemove) MarshalJSON() ([]byte, error) {
	a := (*aliasReplyKeyboardRemove)(m)
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (m *ReplyKeyboardRemove) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var a aliasReplyKeyboardRemove
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	}
	*m = (ReplyKeyboardRemove)(a)
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (m *InlineKeyboardMarkup) MarshalJSON() ([]byte, error) {
	a := (*aliasInlineKeyboardMarkup)(m)
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (m *InlineKeyboardMarkup) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var a aliasInlineKeyboardMarkup
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	}
	*m = (InlineKeyboardMarkup)(a)
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (m *ForceReply) MarshalJSON() ([]byte, error) {
	a := (*aliasForceReply)(m)
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (m *ForceReply) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var a aliasForceReply
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	}
	*m = (ForceReply)(a)
	return nil
}
