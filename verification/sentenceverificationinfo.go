package verification

import (
	"github.com/jamesrwaugh/kyoro/acquisition"
)

// SentenceVerificationInfo contains information needed for
// a SentenceVerififier to do what it does.
type SentenceVerificationInfo struct {
	AcceptedSentences []acquisition.Translation
	Sentence          acquisition.Translation
	SentenceIndex     int
}
