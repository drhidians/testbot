#!/usr/bin/env python

header = '''package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "encoding/json"

'''

keyboard_types = [
    'ReplyKeyboardMarkup',
    'ReplyKeyboardRemove',
    'InlineKeyboardMarkup',
    'ForceReply',
]

alias_types = ['alias' + typ for typ in keyboard_types]

alias_template = '''type {alias_type} {keyboard_type}
'''

methods_template = '''
// MarshalJSON implements json.Marshaler interface.
func (m *{keyboard_type}) MarshalJSON() ([]byte, error) {
    a := (*{alias_type})(m)
    b, err := json.Marshal(a)
    if err != nil {
        return nil, err
    }
    return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (m *{keyboard_type}) UnmarshalJSON(b []byte) error {
    var s string
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    var a {alias_type}
    if err := json.Unmarshal([]byte(s), &a); err != nil {
        return err
    }
    *m = ({keyboard_type})(a)
    return nil
}
'''


def replace(template, replacements):
    s = template
    for key, value in replacements.items():
        s = s.replace(key, value)
    return s


def main():
    with open('types_keyboards.go', 'w') as f:
        f.write(header)
        for typ, alias in zip(keyboard_types, alias_types):
            f.write(replace(alias_template, {
                '{keyboard_type}': typ,
                '{alias_type}': alias,
            }))
        for typ, alias in zip(keyboard_types, alias_types):
            f.write(replace(methods_template, {
                '{keyboard_type}': typ,
                '{alias_type}': alias,
            }))


if __name__ == '__main__':
    main()
