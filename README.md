# emojid

UUID-shaped IDs made of emojis. Format: `8-4-4-4-12` (32 emojis + 4 dashes).

## Install

```bash
go get github.com/pizza-power/emojid@v0.1.0
```

## Example

From the repo root, run the example to print a new emoji ID to the terminal:

```bash
cd example
go run .
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
	_, _ = emojid.Parse(s)
	_ = emojid.Validate("😀😃😄😁😆😅😂🤣-😊😇🙂🙃-😉😌😍🥰-😘😗😙😚-😋😛😝😜🤪🤨🧐🤓😎🥳😤😡")
	_ = emojid.MustParse(s)

	// Custom alphabet
	id, _ = emojid.NewWithAlphabet(emojid.DefaultAlphabet)
	_, _ = emojid.ParseWithAlphabet(s, emojid.DefaultAlphabet)
	
}
```

or to simply print an emoji ID

```go
package main

import (
	"fmt"
	"github.com/pizza-power/emojid"
)

func main() {
	s := emojid.MustNewString()
	fmt.Println(s)
}
```

Errors: `ErrInvalidFormat`, `ErrInvalidToken`, `ErrEntropyFailure`, `ErrAlphabetTooSmall`.
