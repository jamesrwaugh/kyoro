package kyoro

import "log"

type KyoroProduction struct {
}

func (this KyoroProduction) makeSentenceAnkiCard(sentence Translation, options Options) AnkiCard {
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

func (this KyoroProduction) makeKeywordAnkiCard(options Options, inputPhrase string, sentences []Translation) AnkiCard {
	cardFields := map[string]string{
		"japanese": inputPhrase,
	}
	/*if !options.MonoligualMode {
		cardFields["english"] = sentence.English
	}
	for _, sentence := range sentences {
		card := this.makeSentenceAnkiCard(sentence, options)
		anki.AddCard(card)
	}*/
	return AnkiCard{
		DeckName:  options.DeckName,
		ModelName: options.ModelName,
		Fields:    cardFields,
		Tags: []string{
			"kyoro",
		},
	}
}

func (this KyoroProduction) Kyoro(options Options, anki AnkiService, sentenceSource SentenceRetriever, meaningSource MeaningRetriever) bool {
	if !anki.IsConnected() {
		log.Fatal("Could not connect to Anki. Failing.")
		return false
	}
	meanings := meaningSource.GetMeaningforKanji(options.InputPhrase)
	sentences := sentenceSource.GetSentencesforKanji(options.InputPhrase, options.MaxSentences)
	if options.SentencesOnFrontMode {
		for _, sentence := range sentences {
			card := this.makeSentenceAnkiCard(sentence, options)
			anki.AddCard(card)
		}
	} else {

	}

	return true
}
