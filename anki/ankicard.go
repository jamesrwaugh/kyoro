package anki

// AnkiCard represents a card to go into Anki
type AnkiCard struct {
	DeckName  string
	ModelName string
	Fields    map[string]string
	Tags      []string
}
