package kyoro

// AnkiCard represents a card to go into Anki
type AnkiCard struct {
	DeckName  string
	ModelName string
	Fields    map[string]string
	Tags      []string
}

// Sentence has enough information to describe a sentence
type Translation struct {
	Japanese string
	Reading  string
	English  string
}

type MeaningRetriever interface {
	GetMeaningforKanji(kanji string) Translation
}

// SentenceRetriever gets sentences from an outside source for a kanji
type SentenceRetriever interface {
	GetSentencesforKanji(kanji string, maxSentences int) []Translation
}

// AnkiService should add a given card to Anki.
type AnkiService interface {
	IsConnected() bool
	MaxSentencesPerCard() int
	AddCard(card AnkiCard) bool
}

// Options tells Kyoro how to help you.
type Options struct {
	// The Japanese input phrase to generate cards for.
	// This can be a single Kanji, a grammer point, a clause
	// phrase, or anything else.
	InputPhrase string
	// The deck to place new cards in
	DeckName string
	// The card model to submit new notes with.
	// The model must contain the following fieleds:
	//	- "japanese"
	//	- "reading"
	//	- "english"
	ModelName string
	// If true, a card for each found sentence will be created
	//	(up to MaxSentences), where the Japanese sentence is on the front,
	//  and the meanings and readings on the back
	// If false, a card with only the input phrase on the front will be created
	//  and multiple Japanese example sentences with readings using the phrase
	//	input phrase on the back
	SentencesOnFrontMode bool
	// If true, no English reading will be placed on the back of the
	// generated card.
	MonoligualMode bool
	// The maximum number of sentences to pull for the input term.
	MaxSentences int
	// When true, no confirmation dialogs will appear and any sentences
	// pulled will be added as cards unequivocally.
	SilentMode bool
}

// Kyoro is the driver.
type Kyoro interface {
	Kyoro(options Options, anki AnkiService, sentenceSource SentenceRetriever) bool
}
