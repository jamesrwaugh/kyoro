package anki

// AnkiService should add a given card to Anki.
type AnkiService interface {
	IsConnected() bool
	MaxSentencesPerCard() int
	HasMiaCardModel() bool
	AddCard(card AnkiCard) bool
}
