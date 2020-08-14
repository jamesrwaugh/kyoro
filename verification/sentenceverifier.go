package verification

// SentenceVerifier will somehow ask the user that a certain
// sentence is acceptable to add to their deck.
type SentenceVerifier interface {
	UserConfirmSentence(sentence *SentenceVerificationInfo) bool
}
