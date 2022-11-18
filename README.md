# tgrep: Templated Grep

`tgrep` provides basic templating for regular expressions.

## Usage:

`tg`: The `tg` command parses every argument as a templated regular expression and outputs all the templated arguments as regular expressions, one per line.

`tgrep`: The `tgrep` command parses each argument like the `tg` command, but it also reads from `os.Stdin` and outputs every line that matches.

## About Templates

Regex templates allow you to easily create regular expresions that match complex strings without needing to escape any characters.

For example, if you wanted to do a regex search for the string `[code 45]` you would want a regex like `\[code [0-9]+\]`, templing allows that!

### How to use templating

Wrap the portion of your string you **want** to be a regex in double braces. E.g.:

```
[code {{[0-9]+}}]
```

The escaping will be done for you!

### Template shortcuts

There are often common strings that you want to template out, such as email addresses, UUID's, or numbers. You can easily template those too!

```
[code {{number}}]
```

These shortcuts are:

- `number`: a floating point number with optional decimal places (e.g. : 1, 2.0, 3.345)
- `int`: An integer
- `UUID`: A common 8-4-4-4-12 UUID, e.g.: `E24474D7-681E-4293-A450-65F4134E5C36`
