package emojid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// EmojiID is a UUID-shaped identifier composed of emoji tokens.
// Layout: 8-4-4-4-12 emojis (total 32 emojis + 4 dashes).
type EmojiID struct {
	tokens [32]rune
}

// Common errors.
var (
	ErrInvalidFormat    = errors.New("emojid: invalid format")
	ErrInvalidToken     = errors.New("emojid: invalid token (emoji not in alphabet)")
	ErrEntropyFailure   = errors.New("emojid: failed to read crypto randomness")
	ErrAlphabetTooSmall = errors.New("emojid: emoji alphabet must contain at least 2 entries")
)

// DefaultAlphabet is a curated set of single-codepoint emoji.
// Avoids ZWJ sequences, flags, skin tones, and other multi-codepoint grapheme clusters.
var DefaultAlphabet = []rune{
	'😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣',
	'😊', '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰',
	'😘', '😗', '😙', '😚', '😋', '😛', '😝', '😜',
	'🤪', '🤨', '🧐', '🤓', '😎', '🥳', '😤', '😡',
	'🤯', '😱', '😴', '🤤', '😷', '🤒', '🤕', '🤠',
	'😈', '👻', '🤖', '🎃', '🐶', '🐱', '🐭', '🐹',
	'🐰', '🦊', '🐻', '🐼', '🐨', '🐯', '🦁', '🐸',
	'🐵', '🐔', '🐧', '🐦', '🐤', '🐙', '🦑', '🦀',
	'🐠', '🐳', '🦋', '🐞', '🌸', '🌼', '🌻', '🌺',
	'🍎', '🍊', '🍋', '🍉', '🍇', '🍓', '🍒', '🍍',
	'🥑', '🥦', '🥕', '🌶', '🍔', '🍟', '🍕', '🌮',
	'🍣', '🍩', '🍪', '🍫', '🍿', '☕', '🍺', '🍷',
	'⚽', '🏀', '🏈', '⚾', '🎾', '🏐', '🎱', '🏓',
	'🎸', '🎹', '🥁', '🎻', '🎧', '🎮', '🧩', '🎲',
	'🚗', '🚕', '🚌', '🚑', '🚒', '🚜', '✈', '🚀',
	'🛰', '⛵', '🚲', '🛴', '🏠', '🏢', '🏭', '🏰',
	'🌍', '🌙', '⭐', '⚡', '🔥', '💧', '🌈', '❄',
	'💎', '🔒', '🔑', '🧠', '💡', '📦', '🧲', '🧰',
	'🛡', '⚙', '🧪', '🧬', '🔭', '📡', '💾', '🗄',
}

// New returns a new random EmojiID using DefaultAlphabet.
func New() (EmojiID, error) {
	return NewWithAlphabet(DefaultAlphabet)
}

// MustNew is like New but panics on error.
func MustNew() EmojiID {
	id, err := New()
	if err != nil {
		panic(err)
	}
	return id
}

// NewString returns a new random EmojiID as a formatted string.
func NewString() (string, error) {
	id, err := New()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// MustNewString is like NewString but panics on error.
func MustNewString() string {
	id := MustNew()
	return id.String()
}

// NewWithAlphabet returns a new random EmojiID from the provided emoji alphabet.
// Alphabet must contain single-codepoint emoji (runes) and at least 2 entries.
func NewWithAlphabet(alphabet []rune) (EmojiID, error) {
	if len(alphabet) < 2 {
		return EmojiID{}, ErrAlphabetTooSmall
	}

	var id EmojiID

	// We need 32 independent random choices in [0, len(alphabet)).
	// Use rejection sampling from crypto/rand to avoid modulo bias.
	for i := 0; i < len(id.tokens); i++ {
		idx, err := cryptoRandIndex(len(alphabet))
		if err != nil {
			return EmojiID{}, err
		}
		id.tokens[i] = alphabet[idx]
	}

	return id, nil
}

// String formats the EmojiID in the UUID-like layout: 8-4-4-4-12 emojis.
func (e EmojiID) String() string {
	// 32 emojis + 4 dashes.
	var b strings.Builder
	b.Grow(32*utf8.UTFMax + 4)

	writeRunes := func(from, to int) {
		for i := from; i < to; i++ {
			b.WriteRune(e.tokens[i])
		}
	}

	writeRunes(0, 8)
	b.WriteByte('-')
	writeRunes(8, 12)
	b.WriteByte('-')
	writeRunes(12, 16)
	b.WriteByte('-')
	writeRunes(16, 20)
	b.WriteByte('-')
	writeRunes(20, 32)

	return b.String()
}

// Equal compares two EmojiIDs.
func (e EmojiID) Equal(other EmojiID) bool {
	return e.tokens == other.tokens
}

// IsZero reports whether this is the zero value (all tokens are 0 runes).
func (e EmojiID) IsZero() bool {
	var z EmojiID
	return e.tokens == z.tokens
}

// Parse parses an EmojiID string in 8-4-4-4-12 emoji layout using DefaultAlphabet.
func Parse(s string) (EmojiID, error) {
	return ParseWithAlphabet(s, DefaultAlphabet)
}

// MustParse panics if Parse fails.
func MustParse(s string) EmojiID {
	id, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return id
}

// ParseWithAlphabet parses an EmojiID string in 8-4-4-4-12 layout and validates
// that every emoji token is present in the given alphabet.
func ParseWithAlphabet(s string, alphabet []rune) (EmojiID, error) {
	if len(alphabet) < 2 {
		return EmojiID{}, ErrAlphabetTooSmall
	}

	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return EmojiID{}, ErrInvalidFormat
	}

	// Expect emoji counts: 8,4,4,4,12
	want := []int{8, 4, 4, 4, 12}
	total := 0
	for _, w := range want {
		total += w
	}

	var tokens []rune
	tokens = make([]rune, 0, total)

	for i, p := range parts {
		r := []rune(p)
		if len(r) != want[i] {
			return EmojiID{}, ErrInvalidFormat
		}
		tokens = append(tokens, r...)
	}

	if len(tokens) != 32 {
		return EmojiID{}, ErrInvalidFormat
	}

	allowed := make(map[rune]struct{}, len(alphabet))
	for _, r := range alphabet {
		allowed[r] = struct{}{}
	}

	var id EmojiID
	for i := 0; i < 32; i++ {
		if _, ok := allowed[tokens[i]]; !ok {
			return EmojiID{}, fmt.Errorf("%w: %q", ErrInvalidToken, string(tokens[i]))
		}
		id.tokens[i] = tokens[i]
	}

	return id, nil
}

// Validate reports whether s is a valid EmojiID formatted string using DefaultAlphabet.
func Validate(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// Tokens returns the underlying 32 emoji tokens as a slice copy.
func (e EmojiID) Tokens() []rune {
	out := make([]rune, 32)
	copy(out, e.tokens[:])
	return out
}

// --- internal randomness helpers ---

func cryptoRandIndex(n int) (int, error) {
	if n <= 0 {
		return 0, ErrAlphabetTooSmall
	}

	// Rejection sampling using a random byte stream.
	// We’ll draw uint16 values to comfortably cover alphabets up to 65535.
	var buf [2]byte
	max := uint32(1<<16) // 65536
	limit := max - (max % uint32(n))

	for {
		if _, err := rand.Read(buf[:]); err != nil {
			return 0, ErrEntropyFailure
		}
		v := uint32(buf[0])<<8 | uint32(buf[1]) // 0..65535
		if v < limit {
			return int(v % uint32(n)), nil
		}
	}
}
