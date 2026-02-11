# emojid

UUID-shaped IDs made of emojis. Format: `8-4-4-4-12` (32 emojis + 4 dashes).

## Install

```bash
go get github.com/pizza-power/emojid@v0.1.0
```

## Usage

```go
package main

import "github.com/pizza-power/emojid"

func main() {
	// Generate
	id, _ := emojid.New()
	s, _ := emojid.NewString()

	// Or panic on error
	id = emojid.MustNew()
	s = emojid.MustNewString()

	// Format & compare
	_ = id.String()
	_ = id.Equal(emojid.MustParse(s))
	_ = id.IsZero()

	// Parse & validate
	parsed, _ := emojid.Parse(s)
	_ = emojid.Validate("😀😃😄😁😆😅😂🤣-😊😇🙂🙃-😉😌😍🥰-😘😗😙😚-😋😛😝😜🤪🤨🧐🤓😎🥳😤😡")
	_ = emojid.MustParse(s)

	// Custom alphabet
	id, _ = emojid.NewWithAlphabet(emojid.DefaultAlphabet)
	parsed, _ = emojid.ParseWithAlphabet(s, emojid.DefaultAlphabet)
}
```

Errors: `ErrInvalidFormat`, `ErrInvalidToken`, `ErrEntropyFailure`, `ErrAlphabetTooSmall`.
