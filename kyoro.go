package kyoro

import "log"

type KyoroProduction struct {
}

func (this KyoroProduction) makeAnkiCard(options Options, sentence Sentence) AnkiCard {
	// TODO: Sentence-front is not supported.
	cardFields := map[string]string{
		"japanese": sentence.Japanese,
		"reading":  sentence.Reading,
	}
	if !options.MonoligualMode {
		cardFields["english"] = sentence.English
	}
	return AnkiCard{
		DeckName:  options.DeckName,
		ModelName: options.ModelName,
		Fields:    cardFields,
		Tags: []string{
			"kyoro",
		},
	}
}

func (this KyoroProduction) Kyoro(options Options, anki AnkiService, sentenceSource SentenceRetriever) bool {
	if !anki.IsConnected() {
		log.Fatal("Could not connect to Anki. Failing.")
		return false
	}
	sentences := sentenceSource.GetSentencesforKanji(options.InputPhrase, options.MaxSentences)
	for _, sentence := range sentences {
		card := this.makeAnkiCard(options, sentence)
		anki.AddCard(card)
	}
	return true
}
