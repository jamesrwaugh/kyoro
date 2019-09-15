package acquisition

// SentenceRetriever should somehow get a list of sentences with Japananese,
// English, and a reading for a given kanji.
type SentenceRetriever interface {
	GetSentencesforKanji(kanji string, maxSentences int) []Translation
}
