package acquisition

// MeaningRetriever should somehow get a translation of a given kanji.
type MeaningRetriever interface {
	GetMeaningforKanji(kanji string) Translation
}
