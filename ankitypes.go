package kyoro

// AnkiCard represents a card to go into Anki
type AnkiCard struct {
	DeckName  string
	ModelName string
	Fields    map[string]string
	Front     string
	Back      string
	Gags      []string
}

// Sentence has enough information to describe a sentence
type Sentence struct {
	Japanese string
	Reading  string
	English  string
}

// SentenceRetriever gets sentences from an outside source for a kanji
type SentenceRetriever interface {
	GetSentencesforKanji(kanji string) []Sentence
}

// AnkiService should add a given card to Anki.
type AnkiService interface {
	IsConnected() bool
	MaxSentencesPerCard() int
	AddCard(card AnkiCard) bool
}

// Kyoro is the driver.
type Kyoro struct {
	Service        AnkiService
	SentenceGetter SentenceRetriever
}
