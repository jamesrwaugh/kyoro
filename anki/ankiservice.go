package anki

// AnkiService should add a given card to Anki.
type AnkiService interface {
	IsConnected() bool
	HasMiaCardModel() bool
	AddCard(card AnkiCard) bool
}
