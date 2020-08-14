package verification

import (
	"github.com/jamesrwaugh/kyoro/acquisition"
)

// SentenceVerifier will somehow ask the user that a certain
// sentence is acceptable to add to their deck.
type SentenceVerifier interface {
	UserConfirmSentence(sentence *acquisition.Translation) bool
}
