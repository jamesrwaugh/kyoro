package anki

// AnkiService should add a given card to Anki.
type AnkiService interface {
	IsConnected() bool
	MaxSentencesPerCard() int
	AddCard(card AnkiCard) bool
}
