#!/usr/bin/env python

header = '''package telegram

// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT CHANGE IT

import "context"
'''

methods = '''
AnswerCallbackQuery CallbackQueryAnswer
DeleteMessage       DeletedMessage
'''.split()

method_template = '''
// https://core.telegram.org/bots/api#{api_method_lower}
func (b *bot) {method}(ctx context.Context, v *{value}) error {
    var ok bool
    if err := b.do(ctx, "{api_method}", v, &ok); err != nil {
        return err
    }
    if !ok {
        return ErrNotAnswered
    }
    return nil
}
'''


def replace(template, replacements):
    s = template
    for key, value in replacements.items():
        s = s.replace(key, value)
    return s


def main():
    pairs = zip(methods[::2], methods[1::2])

    with open('methods_bool.go', 'w') as f:
        f.write(header)
        for method, value in pairs:
            api_method = method[0].lower() + method[1:]

            f.write(replace(method_template, {
                '{method}': method,
                '{value}': value,
                '{api_method}': api_method,
                '{api_method_lower}': api_method.lower(),
            }))


if __name__ == '__main__':
    main()
