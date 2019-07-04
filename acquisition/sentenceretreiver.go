package acquisition

type SentenceRetriever interface {
	GetSentencesforKanji(kanji string, maxSentences int) []Translation
}
