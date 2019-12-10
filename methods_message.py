#!/usr/bin/env python

header = '''package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "context"
'''

methods = '''
SendMessage            TextMessage
ForwardMessage         ForwardedMessage
SendPhoto              PhotoMessage
SendAudio              AudioMessage
SendDocument           DocumentMessage
SendSticker            StickerMessage
SendVideo              VideoMessage
SendVoice              VoiceMessage
SendVoiceNote          VoiceNoteMessage
SendLocation           LocationMessage
SendVenue              VenueMessage
SendContact            ContactMessage
EditMessageText        MessageText
EditMessageCaption     MessageCaption
EditMessageReplyMarkup MessageReplyMarkup
'''.split()

not_implemented_messages = '''
VideoMessage
VoiceMessage
VoiceNoteMessage
LocationMessage
VenueMessage
ContactMessage
'''.split()

method_template = '''
// https://core.telegram.org/bots/api#{api_method_lower}
func (b *bot) {method}(ctx context.Context, m *{message}) (*Message, error) {
    var v *Message
    if err := b.do(ctx, "{api_method}", m, &v); err != nil {
        return nil, err
    }
    return v, nil
}
'''


def replace(template, replacements):
    s = template
    for key, value in replacements.items():
        s = s.replace(key, value)
    return s


def main():
    pairs = zip(methods[::2], methods[1::2])

    with open('methods_message.go', 'w') as f:
        f.write(header)
        for method, message in pairs:
            if message in not_implemented_messages:
                continue

            api_method = method[0].lower() + method[1:]

            f.write(replace(method_template, {
                '{method}': method,
                '{message}': message,
                '{api_method}': api_method,
                '{api_method_lower}': api_method.lower(),
            }))


if __name__ == '__main__':
    main()
